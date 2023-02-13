module github.com/giantswarm/valuemodifier

go 1.15

require (
	github.com/ProtonMail/go-crypto v0.0.0-20220714114130-e85cedf506cd
	github.com/cenkalti/backoff/v3 v3.2.2 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/giantswarm/microerror v0.4.0
	github.com/hashicorp/go-hclog v1.2.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.1 // indirect
	github.com/hashicorp/vault/api v1.9.0
	github.com/spf13/cast v1.5.0
	golang.org/x/time v0.0.0-20220609170525-579cf78fd858 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/go-ldap/ldap/v3 v3.1.10 => github.com/go-ldap/ldap/v3 v3.4.3
	github.com/gogo/protobuf v1.1.1 => github.com/gogo/protobuf v1.3.2
	github.com/prometheus/client_golang v1.11.0 => github.com/prometheus/client_golang v1.12.2
	github.com/prometheus/client_golang v1.4.0 => github.com/prometheus/client_golang v1.12.2
)
