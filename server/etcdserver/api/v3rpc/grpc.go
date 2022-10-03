// Copyright 2016 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v3rpc

import (
	"crypto/tls"
	"math"

	pb "github.com/lfch/etcd-io/api/v3/etcdserverpb"
	"github.com/lfch/etcd-io/client/v3/credentials"
	"github.com/lfch/etcd-io/server/v3/etcdserver"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/lfch/grpc"
	"github.com/lfch/grpc/health"
	healthpb "github.com/lfch/grpc/health/grpc_health_v1"
	"go.opentelemetry.io/contrib/instrumentation/github.com/lfch/grpc/otelgrpc"
)

const (
	grpcOverheadBytes = 512 * 1024
	maxSendBytes      = math.MaxInt32
)

func Server(s *etcdserver.EtcdServer, tls *tls.Config, interceptor grpc.UnaryServerInterceptor, gopts ...grpc.ServerOption) *grpc.Server {
	var opts []grpc.ServerOption
	opts = append(opts, grpc.CustomCodec(&codec{}))
	if tls != nil {
		bundle := credentials.NewBundle(credentials.Config{TLSConfig: tls})
		opts = append(opts, grpc.Creds(bundle.TransportCredentials()))
	}
	chainUnaryInterceptors := []grpc.UnaryServerInterceptor{
		newLogUnaryInterceptor(s),
		newUnaryInterceptor(s),
		grpc_prometheus.UnaryServerInterceptor,
	}
	if interceptor != nil {
		chainUnaryInterceptors = append(chainUnaryInterceptors, interceptor)
	}

	chainStreamInterceptors := []grpc.StreamServerInterceptor{
		newStreamInterceptor(s),
		grpc_prometheus.StreamServerInterceptor,
	}

	if s.Cfg.ExperimentalEnableDistributedTracing {
		chainUnaryInterceptors = append(chainUnaryInterceptors, otelgrpc.UnaryServerInterceptor(s.Cfg.ExperimentalTracerOptions...))
		chainStreamInterceptors = append(chainStreamInterceptors, otelgrpc.StreamServerInterceptor(s.Cfg.ExperimentalTracerOptions...))

	}

	opts = append(opts, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(chainUnaryInterceptors...)))
	opts = append(opts, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(chainStreamInterceptors...)))

	opts = append(opts, grpc.MaxRecvMsgSize(int(s.Cfg.MaxRequestBytes+grpcOverheadBytes)))
	opts = append(opts, grpc.MaxSendMsgSize(maxSendBytes))
	opts = append(opts, grpc.MaxConcurrentStreams(s.Cfg.MaxConcurrentStreams))

	grpcServer := grpc.NewServer(append(opts, gopts...)...)

	pb.RegisterKVServer(grpcServer, NewQuotaKVServer(s))
	pb.RegisterWatchServer(grpcServer, NewWatchServer(s))
	pb.RegisterLeaseServer(grpcServer, NewQuotaLeaseServer(s))
	pb.RegisterClusterServer(grpcServer, NewClusterServer(s))
	pb.RegisterAuthServer(grpcServer, NewAuthServer(s))
	pb.RegisterMaintenanceServer(grpcServer, NewMaintenanceServer(s))

	// server should register all the services manually
	// use empty service name for all etcd services' health status,
	// see https://github.com/grpc/grpc/blob/master/doc/health-checking.md for more
	hsrv := health.NewServer()
	hsrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(grpcServer, hsrv)

	// set zero values for metrics registered for this grpc server
	grpc_prometheus.Register(grpcServer)

	return grpcServer
}
