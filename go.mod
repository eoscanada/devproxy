module github.com/eoscanada/devproxy

go 1.12

require (
	contrib.go.opencensus.io/exporter/stackdriver v0.12.6
	github.com/dfuse-io/derr v0.0.0-20200406214256-c690655246a1
	github.com/dfuse-io/dfuse-saas-priv v0.0.0-20200416135030-e98fed343fc8
	github.com/dfuse-io/logging v0.0.0-20200407175011-14021b7a79af
	github.com/gogo/protobuf v1.2.1
	github.com/golang/protobuf v1.4.0
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190717153623-606c73359dba
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/mwitkow/grpc-proxy v0.0.0-20181017164139-0f1106ef9c76
	github.com/prometheus/client_golang v1.1.0 // indirect
	go.opencensus.io v0.22.2
	go.uber.org/zap v1.14.0
	google.golang.org/grpc v1.28.1
)

// This is required to fix build where 0.1.0 version is not considered a valid version because a v0 line does not exists
// We replace with same commit, simply tricking go and tell him that's it's actually version 0.0.3
replace github.com/census-instrumentation/opencensus-proto v0.1.0-0.20181214143942-ba49f56771b8 => github.com/census-instrumentation/opencensus-proto v0.0.3-0.20181214143942-ba49f56771b8
