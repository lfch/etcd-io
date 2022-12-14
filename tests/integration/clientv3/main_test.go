// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package clientv3test

import (
	"testing"

	"github.com/lfch/etcd-io/client/pkg/v3/testutil"
)

func TestMain(m *testing.M) {
	testutil.MustTestMainWithLeakDetection(m)
}
