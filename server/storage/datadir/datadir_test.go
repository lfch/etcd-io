package datadir_test

import (
	"testing"

	"github.com/lfch/etcd-io/server/v3/storage/datadir"
	"github.com/stretchr/testify/assert"
)

func TestToBackendFileName(t *testing.T) {
	result := datadir.ToBackendFileName("/dir/data-dir")
	assert.Equal(t, "/dir/data-dir/member/snap/db", result)
}

func TestToMemberDir(t *testing.T) {
	result := datadir.ToMemberDir("/dir/data-dir")
	assert.Equal(t, "/dir/data-dir/member", result)
}

func TestToSnapDir(t *testing.T) {
	result := datadir.ToSnapDir("/dir/data-dir")
	assert.Equal(t, "/dir/data-dir/member/snap", result)
}

func TestToWalDir(t *testing.T) {
	result := datadir.ToWalDir("/dir/data-dir")
	assert.Equal(t, "/dir/data-dir/member/wal", result)
}

func TestToWalDirSlash(t *testing.T) {
	result := datadir.ToWalDir("/dir/data-dir/")
	assert.Equal(t, "/dir/data-dir/member/wal", result)
}
