package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	stackdriverPropagation "contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"go.opencensus.io/plugin/ochttp"
)

func NewReverseProxy(target *url.URL, stripQuerystring bool) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		if stripQuerystring {
			req.URL.RawQuery = ""
		}
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
		req.Header.Set("Host", target.Host)
		if _, ok := req.Header["User-Agent"]; !ok {
			// Explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}

	return &httputil.ReverseProxy{
		Director: director,
		ModifyResponse: func(response *http.Response) error {
			response.Header.Del("X-Trace-ID")
			return nil
		},
		Transport: &ochttp.Transport{
			Propagation: &stackdriverPropagation.HTTPFormat{},
		},
	}
}
