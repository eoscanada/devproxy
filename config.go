package main

import (
	pbreflect "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

type config struct {
	serviceToEndpoint    map[string]string
	allServices          []string
	fileContainingSymbol map[string]*pbreflect.ServerReflectionResponse
	extensionNumbers     map[string]*pbreflect.ServerReflectionResponse
	filesByFilename       map[string]*pbreflect.ServerReflectionResponse
}

func newConfig() *config {
	return &config{
		serviceToEndpoint:    map[string]string{},
		fileContainingSymbol: map[string]*pbreflect.ServerReflectionResponse{},
		extensionNumbers:     map[string]*pbreflect.ServerReflectionResponse{},
		filesByFilename:       map[string]*pbreflect.ServerReflectionResponse{},
	}
}
