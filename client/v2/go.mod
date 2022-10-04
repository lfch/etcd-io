module github.com/lfch/etcd-io/client/v2

go 1.19

require (
	github.com/lfch/etcd-io/api/v3 v3.6.0-alpha.0
	github.com/lfch/etcd-io/client/pkg/v3 v3.6.0-alpha.0
	sigs.k8s.io/json v0.0.0-20211020170558-c049b76a60c6
)

require (
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/testify v1.7.2 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/lfch/etcd-io/api/v3 => ../../api
	github.com/lfch/etcd-io/client/pkg/v3 => ../pkg
)

// Bad imports are sometimes causing attempts to pull that code.
// This makes the error more explicit.
replace (
	github.com/lfch/etcd-io => ./FORBIDDEN_DEPENDENCY
	github.com/lfch/etcd-io/pkg/v3 => ./FORBIDDED_DEPENDENCY
	github.com/lfch/etcd-io/tests/v3 => ./FORBIDDEN_DEPENDENCY
	github.com/lfch/etcd-io/v3 => ./FORBIDDEN_DEPENDENCY
)
