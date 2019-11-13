module github.com/eoscanada/devproxy

go 1.12

require (
	cloud.google.com/go v0.43.0
	contrib.go.opencensus.io/exporter/stackdriver v0.12.6
	github.com/eoscanada/bstream v1.6.2
	github.com/eoscanada/dbilling v1.3.1
	github.com/eoscanada/derr v0.3.9
	github.com/eoscanada/dstore v0.0.1
	github.com/eoscanada/logging v0.6.5
	github.com/eoscanada/search v0.0.0-20190822144146-0bce109a3bd8
	github.com/gogo/protobuf v1.2.1
	github.com/golang/protobuf v1.3.2
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190717153623-606c73359dba
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/karixtech/zapdriver v1.1.7-0.20190304072941-7b5d38c10286 // indirect
	github.com/mwitkow/grpc-proxy v0.0.0-20181017164139-0f1106ef9c76
	github.com/prometheus/client_golang v1.1.0 // indirect
	github.com/stretchr/testify v1.4.0 // indirect
	go.opencensus.io v0.22.1
	go.uber.org/zap v1.10.0
	google.golang.org/grpc v1.23.0
)

// This is required to fix build where 0.1.0 version is not considered a valid version because a v0 line does not exists
// We replace with same commit, simply tricking go and tell him that's it's actually version 0.0.3
replace github.com/census-instrumentation/opencensus-proto v0.1.0-0.20181214143942-ba49f56771b8 => github.com/census-instrumentation/opencensus-proto v0.0.3-0.20181214143942-ba49f56771b8
