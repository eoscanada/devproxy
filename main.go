package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/dfuse-io/derr"
	pbdevproxy "github.com/dfuse-io/dfuse-saas-priv/pb/dfuse/devproxy/v1"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	proxy "github.com/mwitkow/grpc-proxy/proxy"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	pbreflect "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

var flagHTTPAddr = flag.String("http-addr", ":8001", "HTTP proxy listen address")
var flagListenAddr = flag.String("listen-addr", ":9001", "gRPC listen address")
var flagProxies = flag.String("proxies", "localhost:8080/v0/state,localhost:9999/v1/chain", "Comma-separated list of service:port/path to reverse proxy through HTTP.")
var flagServices = flag.String("services", "localhost:7001*,localhost:7002,localhost:7003,localhost:7004", "Comma-separated list of service:port to reverse proxy and cumulate Reflection endpoints.")

func main() {
	flag.Parse()
	setupLogger()

	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)

	go LaunchGRPCServer(waitGroup)
	go LaunchHTTPProxy(waitGroup)

	zlog.Info("ready")
	waitGroup.Wait()
}

func LaunchGRPCServer(waitGroup *sync.WaitGroup) {
	services := strings.SplitN(*flagServices, ",", -1)
	conf := newConfig()

	// Aggregate all the methods supported
	derr.ErrorCheck("discover", discover(services, conf))

	lis, err := net.Listen("tcp", *flagListenAddr)
	if err != nil {
		zlog.Fatal("failed listening grpc", zap.String("grpc_listen_addr", *flagListenAddr), zap.Error(err))
	}

	zlog.Debug("known services", zap.Strings("services", conf.allServices))
	srv := &ReflectServer{conf: conf}

	gs := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_prometheus.UnaryServerInterceptor,
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_prometheus.StreamServerInterceptor,
		),
		grpc.CustomCodec(proxy.Codec()),
		grpc.UnknownServiceHandler(proxy.TransparentHandler(srv.Director)),
	)

	pbdevproxy.RegisterDevproxyServer(gs, srv)
	pbreflect.RegisterServerReflectionServer(gs, srv)

	// When queries on the reflection endpoint, return a UNION of all the services behind/below.
	// Reverse proxy any incoming queries to the right backend service.
	zlog.Info("listening & serving gRPC content", zap.String("grpc_listen_addr", *flagListenAddr))

	if err := gs.Serve(lis); err != nil {
		waitGroup.Done()
		zlog.Fatal("error on gRPC proxy", zap.Error(err))
	}
}

type ProxyInfo struct {
	Scheme string
	Host   string
	Path   string
}

func (p *ProxyInfo) URL() *url.URL {
	return &url.URL{Scheme: p.Scheme, Host: p.Host}
}

func (p *ProxyInfo) String() string {
	return p.URL().String()
}

func LaunchHTTPProxy(waitGroup *sync.WaitGroup) {
	proxies := strings.Split(*flagProxies, ",")
	regex := *regexp.MustCompile(`([a-z0-9-_\.]+:[0-9]{2,5})(/[^,]+)`)
	matches := regex.FindAllStringSubmatch(*flagProxies, -1)

	if len(matches) != len(proxies) {
		zlog.Fatal(fmt.Sprintf("Flag --proxies is invalid, check the formatting, was able to only parse %d out of %d", len(matches), len(proxies)))
	}

	proxiesInfo := make([]*ProxyInfo, len(matches))
	for i, match := range matches {
		if len(match) != 3 {
			zlog.Fatal(fmt.Sprintf("Flag --proxies value %q is invalid, check the its formatting", proxies[i]))
		}

		proxiesInfo[i] = &ProxyInfo{Scheme: "http", Host: match[1], Path: match[2]}
	}

	router := mux.NewRouter()
	for _, proxyInfo := range proxiesInfo {
		zlog.Info("proxy info", zap.Stringer("info", proxyInfo))
		router.PathPrefix(proxyInfo.Path).Handler(NewReverseProxy(proxyInfo.URL(), false))
	}

	zlog.Info("serving HTTP proxy", zap.String("http_addr", *flagHTTPAddr))
	server := &http.Server{Addr: *flagHTTPAddr, Handler: handlers.CompressHandlerLevel(router, gzip.BestSpeed)}
	err := server.ListenAndServe()
	if err != nil {
		waitGroup.Done()
		zlog.Fatal("error on HTTP proxy", zap.Error(err))
	}
}

// func mustToInt(input string) int {
// 	value, err := strconv.ParseInt(input, 10)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return vlaue
// }

// chainRouter.PathPrefix("/v1/chain").Handler(txPushRouter)
