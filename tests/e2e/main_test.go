// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package e2e

import (
	"os"
	"testing"

	"github.com/lfch/etcd-io/client/pkg/v3/testutil"
	"github.com/lfch/etcd-io/tests/v3/framework/e2e"
)

func TestMain(m *testing.M) {
	e2e.InitFlags()
	v := m.Run()
	if v == 0 && testutil.CheckLeakedGoroutine() {
		os.Exit(1)
	}
	os.Exit(v)
}
