module github.com/hashicorp/vault/api

go 1.13

replace github.com/hashicorp/vault/sdk => ../sdk

require (
	github.com/cenkalti/backoff/v3 v3.2.2
	github.com/go-test/deep v1.0.8
	github.com/hashicorp/errwrap v1.1.0
	github.com/hashicorp/go-cleanhttp v0.5.2
	github.com/hashicorp/go-hclog v1.1.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-retryablehttp v0.7.0
	github.com/hashicorp/go-rootcerts v1.0.2
	github.com/hashicorp/go-secure-stdlib/parseutil v0.1.3
	github.com/hashicorp/hcl v1.0.1-vault-3
	github.com/hashicorp/vault v1.9.4
	github.com/hashicorp/vault/sdk v0.4.1
	github.com/mitchellh/mapstructure v1.4.3
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd
	golang.org/x/time v0.0.0-20211116232009-f0f3c7e86c11
	gopkg.in/square/go-jose.v2 v2.6.0
)
