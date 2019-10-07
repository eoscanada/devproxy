package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/eoscanada/derr"
	"github.com/gogo/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pbreflect "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

func dialOptions(endpoint string) (opts []grpc.DialOption) {
	if strings.Contains(endpoint, "*") {
		zlog.Info("with transport credentials")
		creds := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		})
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		zlog.Info("insecure endpoint")
		opts = append(opts, grpc.WithInsecure())
	}
	return opts
}

func discover(services []string, conf *config) error {
	filesRequested := map[string]bool{}

	for _, srv := range services {
		srv = strings.TrimSpace(srv)
		target := strings.Replace(srv, "*", "", -1)

		zlog.Info("querying service " + srv)
		// CALL the reflection API there
		opts := dialOptions(srv)

		conn, err := grpc.Dial(target, opts...)
		derr.ErrorCheck("dialing to service "+srv, err)

		client := pbreflect.NewServerReflectionClient(conn)
		stream, err := client.ServerReflectionInfo(context.Background())
		derr.ErrorCheck("setting up client", err)

		err = stream.Send(&pbreflect.ServerReflectionRequest{
			Host:           target,
			MessageRequest: &pbreflect.ServerReflectionRequest_ListServices{ListServices: "*"},
		})
		derr.ErrorCheck("sending list services request", err)

		resp, err := stream.Recv()
		derr.ErrorCheck("receiving list services response", err)

		zlog.Info("reflection list services response", zap.Any("response", toMap(resp)))

		switch msg := resp.MessageResponse.(type) {
		case *pbreflect.ServerReflectionResponse_ListServicesResponse:
			var reqs int
			for _, serviceResp := range msg.ListServicesResponse.Service {
				serviceName := serviceResp.Name
				conf.allServices = append(conf.allServices, serviceName)
				conf.serviceToEndpoint[serviceName] = srv

				err = stream.Send(&pbreflect.ServerReflectionRequest{
					Host:           target,
					MessageRequest: &pbreflect.ServerReflectionRequest_FileContainingSymbol{FileContainingSymbol: serviceName},
				})
				derr.ErrorCheck("sending reflection request", err)

				err = stream.Send(&pbreflect.ServerReflectionRequest{
					Host:           target,
					MessageRequest: &pbreflect.ServerReflectionRequest_AllExtensionNumbersOfType{AllExtensionNumbersOfType: serviceName},
				})
				derr.ErrorCheck("sending reflection request", err)

				reqs += 2
			}

			for i := 0; i < reqs; i++ {
				resp, err := stream.Recv()
				derr.ErrorCheck("receiving reflection response", err)

				origReq := resp.OriginalRequest
				zlog.Debug("reflection request response", zap.Any("request", toMap(origReq)), zap.Any("resonse", toMap(resp)))

				switch msg := resp.MessageResponse.(type) {
				case *pbreflect.ServerReflectionResponse_AllExtensionNumbersResponse:
					r := msg.AllExtensionNumbersResponse
					origSymbol := origReq.MessageRequest.(*pbreflect.ServerReflectionRequest_AllExtensionNumbersOfType).AllExtensionNumbersOfType
					zlog.Debug("all extensions number",
						zap.String("base_type", r.BaseTypeName),
						zap.Int32s("extension_number", r.ExtensionNumber),
						zap.String("original_symbol", origSymbol),
					)

					conf.extensionNumbers[origSymbol] = resp

				case *pbreflect.ServerReflectionResponse_FileDescriptorResponse:
					r := msg.FileDescriptorResponse

					var filenames []string
					for _, descFile := range r.FileDescriptorProto {
						desc := &descriptor.FileDescriptorProto{}
						err = proto.Unmarshal(descFile, desc)
						derr.ErrorCheck("unmarshal file descriptor proto", err)

						filenames = append(filenames, desc.Dependency...)
					}

					for _, fileName := range filenames {
						if !filesRequested[fileName] {
							err = stream.Send(&pbreflect.ServerReflectionRequest{
								Host:           target,
								MessageRequest: &pbreflect.ServerReflectionRequest_FileByFilename{FileByFilename: fileName},
							})
							reqs++

							filesRequested[fileName] = true
						}
					}

					switch origPayload := origReq.MessageRequest.(type) {
					case *pbreflect.ServerReflectionRequest_FileContainingSymbol:
						origSymbol := origPayload.FileContainingSymbol
						conf.fileContainingSymbol[origSymbol] = resp

					case *pbreflect.ServerReflectionRequest_FileByFilename:
						origFilename := origPayload.FileByFilename
						conf.filesByFilename[origFilename] = resp
					}

				case *pbreflect.ServerReflectionResponse_ErrorResponse:
					zlog.Warn("received reflection error response", zap.Any("response", toMap(msg.ErrorResponse)))
				default:
					zlog.Warn("an unpextec response type was received but not handled")
				}
			}
		default:
			derr.ErrorCheck("wuut, invalid response to the request we made", fmt.Errorf("we received type %T %+v", msg, msg))
		}

		derr.ErrorCheck("close send", stream.CloseSend())
	}

	return nil
}

func toMap(any interface{}) map[string]interface{} {
	cnt, err := json.Marshal(any)
	derr.ErrorCheck("marshal response", err)

	out := map[string]interface{}{}
	derr.ErrorCheck("unmarshal response", json.Unmarshal(cnt, &out))

	return out
}
