package main

import (
	"context"
	"fmt"
	"io"
	"strings"

	pbdevproxy "github.com/eoscanada/devproxy/pb/dfuse/devproxy/v1"
	proxy "github.com/mwitkow/grpc-proxy/proxy"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	pbreflect "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

type ReflectServer struct {
	conf *config
}

func (s *ReflectServer) Director(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
	var endpoint string
	// check in the headers if `server` is specified, in which case, forward to it directly
	md, ok := metadata.FromIncomingContext(ctx)
	if ok && len(md.Get("server")) != 0 {
		endpoint = md.Get("server")[0]
	} else {
		zlog.Info("full method name", zap.String("method_name", fullMethodName))
		parts := strings.Split(fullMethodName, "/")
		endpoint = s.conf.serviceToEndpoint[parts[1]]
	}

	if endpoint == "" {
		return nil, nil, grpc.Errorf(codes.Unimplemented, "Unknown method or endpoint to reach")
	}

	opts := dialOptions(endpoint)
	opts = append(opts, grpc.WithCodec(proxy.Codec()))

	conn, err := grpc.DialContext(ctx, endpoint, opts...)
	return ctx, conn, err
}

func (s *ReflectServer) ListServers(ctx context.Context, req *pbdevproxy.ListRequest) (*pbdevproxy.ListResponse, error) {
	return &pbdevproxy.ListResponse{
		Servers: s.conf.allServices,
	}, nil
}

func (s *ReflectServer) ServerReflectionInfo(stream pbreflect.ServerReflection_ServerReflectionInfoServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			zlog.Error("error reading message", zap.Error(err))
			break
		}

		switch req := msg.MessageRequest.(type) {
		case *pbreflect.ServerReflectionRequest_FileByFilename:
			err = stream.Send(s.conf.filesByFilename[req.FileByFilename])
		case *pbreflect.ServerReflectionRequest_FileContainingSymbol:
			err = stream.Send(s.conf.fileContainingSymbol[req.FileContainingSymbol])
		case *pbreflect.ServerReflectionRequest_FileContainingExtension:
			err = fmt.Errorf("unimplemented")
		case *pbreflect.ServerReflectionRequest_AllExtensionNumbersOfType:
			err = stream.Send(s.conf.extensionNumbers[req.AllExtensionNumbersOfType])
		case *pbreflect.ServerReflectionRequest_ListServices:
			var services []*pbreflect.ServiceResponse
			seen := map[string]bool{}
			for _, svc := range s.conf.allServices {
				if seen[svc] {
					continue
				}
				seen[svc] = true
				services = append(services, &pbreflect.ServiceResponse{
					Name: svc,
				})
			}

			err = stream.Send(&pbreflect.ServerReflectionResponse{
				ValidHost:       msg.Host, //// wuut anyway?
				OriginalRequest: msg,
				MessageResponse: &pbreflect.ServerReflectionResponse_ListServicesResponse{
					ListServicesResponse: &pbreflect.ListServiceResponse{
						Service: services,
					},
				},
			})
		}
		if err != nil {
			zlog.Error("error sending", zap.Error(err))
			break
		}
	}
	// Take in the command, fetch the right responses from the last global `config`
	return nil
}
