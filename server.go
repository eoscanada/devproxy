package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	proxy "github.com/mwitkow/grpc-proxy/proxy"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	pbreflect "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

type ReflectServer struct {
	conf *config
}

func (s *ReflectServer) Director(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
	fmt.Println("full method name", fullMethodName)
	parts := strings.Split(fullMethodName, "/")
	endpoint := s.conf.serviceToEndpoint[parts[1]]
	if endpoint == "" {
		return nil, nil, grpc.Errorf(codes.Unimplemented, "Unknown method")
	}

	conn, err := grpc.DialContext(ctx, endpoint, grpc.WithInsecure(), grpc.WithCodec(proxy.Codec()))
	return ctx, conn, err

	// md, ok := metadata.FromContext(ctx)
	// if ok {
	// 	// Decide on which backend to dial
	// 	if val, exists := md[":authority"]; exists && val[0] == "staging.api.example.com" {
	// 		// Make sure we use DialContext so the dialing can be cancelled/time out together with the context.
	// 		return grpc.DialContext(ctx, "api-service.staging.svc.local", grpc.WithCodec(proxy.Codec()))
	// 	} else if val, exists := md[":authority"]; exists && val[0] == "api.example.com" {
	// 		return grpc.DialContext(ctx, "api-service.prod.svc.local", grpc.WithCodec(proxy.Codec()))
	// 	}
	// }

	// return grpc.DialContext(
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

		cnt, _ := json.MarshalIndent(msg, "", "  ")
		fmt.Println("RRRRRRRRRITA", string(cnt))

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
			fmt.Println("MAMA", s.conf.allServices)

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
