// Copyright 2017 The etcd Authors
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

package clientv3_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/lfch/etcd-io/client/pkg/v3/testutil"
	integration2 "github.com/lfch/etcd-io/tests/v3/framework/integration"
	"github.com/lfch/etcd-io/tests/v3/integration"
)

const (
	dialTimeout    = 5 * time.Second
	requestTimeout = 10 * time.Second
)

var lazyCluster = integration.NewLazyClusterWithConfig(
	integration2.ClusterConfig{
		Size:                        3,
		WatchProgressNotifyInterval: 200 * time.Millisecond,
		DisableStrictReconfigCheck:  true})

func exampleEndpoints() []string { return lazyCluster.EndpointsV3() }

func forUnitTestsRunInMockedContext(_ func(), example func()) {
	// For integration tests runs in the provided environment
	example()
}

// TestMain sets up an etcd cluster if running the examples.
func TestMain(m *testing.M) {
	testutil.ExitInShortMode("Skipping: the tests require real cluster")

	tempDir, err := os.MkdirTemp(os.TempDir(), "etcd-integration")
	if err != nil {
		log.Printf("Failed to obtain tempDir: %v", tempDir)
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir)

	err = os.Chdir(tempDir)
	if err != nil {
		log.Printf("Failed to change working dir to: %s: %v", tempDir, err)
		os.Exit(1)
	}
	log.Printf("Running tests (examples) in dir(%v): ...", tempDir)
	v := m.Run()
	lazyCluster.Terminate()

	if v == 0 {
		testutil.MustCheckLeakedGoroutine()
	}
	os.Exit(v)
}
