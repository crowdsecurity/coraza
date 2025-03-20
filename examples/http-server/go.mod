module github.com/crowdsecurity/coraza/v3/examples/http-server

go 1.23.0

require github.com/corazawaf/coraza/v3 v3.2.1

require (
	github.com/corazawaf/libinjection-go v0.2.2 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/magefile/mage v1.15.1-0.20241126214340-bdc92f694516 // indirect
	github.com/petar-dambovaliev/aho-corasick v0.0.0-20240411101913-e07a1f0e8eb4 // indirect
	github.com/tidwall/gjson v1.18.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/valllabh/ocsf-schema-golang v1.0.3 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	rsc.io/binaryregexp v0.2.0 // indirect
)

replace github.com/corazawaf/coraza/v3 => github.com/crowdsecurity/coraza/v3 v3.0.0-20250320223126-d7f736e6301f
