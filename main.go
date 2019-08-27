package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/eoscanada/logging"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	proxy "github.com/mwitkow/grpc-proxy/proxy"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	pbreflect "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

var flagListenAddr = flag.String("listen-addr", ":9001", "gRPC listen address")
var flagServices = flag.String("services", "", "Comma-separated list of service:port to reverse proxy and cumulate Reflection endpoints.")

func main() {
	flag.Parse()
	setupLogger()
	//  devproxy --service asljldkjf --listen-addr=:9000 --config /etc/config/map.json

	// /dfuse.devproxy.v1.Infos/GetInfos
	// {"blocks_url": "gs://therightblocks-url-from-the-operator",
	//  "search_version": "aslfdkjdsa"},

	services := strings.SplitN(*flagServices, ",", -1)

	conf := newConfig()

	errorCheck("discover", discover(services, conf))

	zlog.Info("ready")
	// Aggregate all the methods supported

	lis, err := net.Listen("tcp", *flagListenAddr)
	if err != nil {
		zlog.Fatal("failed listening grpc", zap.String("grpc_listen_addr", *flagListenAddr), zap.Error(err))
	}

	srv := &ReflectServer{conf: conf}

	unaryLog, streamLog := logging.ServerInterceptors()
	gs := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_prometheus.UnaryServerInterceptor,
			unaryLog,
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_prometheus.StreamServerInterceptor,
			streamLog,
		),
		grpc.CustomCodec(proxy.Codec()),
		grpc.UnknownServiceHandler(proxy.TransparentHandler(srv.Director)),
	)

	pbreflect.RegisterServerReflectionServer(gs, srv)

	// When queries on the reflection endpoint, return a UNION of all the services behind/below.
	// Reverse proxy any incoming queries to the right backend service.
	zlog.Info("listening & serving gRPC content", zap.String("grpc_listen_addr", *flagListenAddr))
	if err := gs.Serve(lis); err != nil {
		zlog.Fatal("error on gs.Serve", zap.Error(err))
	}
}

func errorCheck(prefix string, err error) {
	if err != nil {
		fmt.Printf("%s: %s\n", prefix, err)
		os.Exit(1)
	}
}
