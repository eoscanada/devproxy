package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strings"

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

		zl := zlog.With(zap.String("service", srv))

		zl.Info("querying service")
		// CALL the reflection API there
		opts := dialOptions(srv)

		conn, err := grpc.Dial(target, opts...)
		errorCheck(zl, "dialing to service error", err)

		client := pbreflect.NewServerReflectionClient(conn)
		stream, err := client.ServerReflectionInfo(context.Background())
		errorCheck(zl, "setting up client error", err)

		err = stream.Send(&pbreflect.ServerReflectionRequest{
			Host:           target,
			MessageRequest: &pbreflect.ServerReflectionRequest_ListServices{ListServices: "*"},
		})
		errorCheck(zl, "sending list services request error", err)

		resp, err := stream.Recv()
		errorCheck(zl, "receiving list services response error", err)

		zl.Info("reflection list services response", zap.Any("response", toMap(resp)))

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
				errorCheck(zl, "sending reflection request error", err)

				err = stream.Send(&pbreflect.ServerReflectionRequest{
					Host:           target,
					MessageRequest: &pbreflect.ServerReflectionRequest_AllExtensionNumbersOfType{AllExtensionNumbersOfType: serviceName},
				})
				errorCheck(zl, "sending reflection request error", err)

				reqs += 2
			}

			for i := 0; i < reqs; i++ {
				resp, err := stream.Recv()
				errorCheck(zl, "receiving reflection response error", err)

				origReq := resp.OriginalRequest
				zl.Info("reflection request response", zap.Any("request", toMap(origReq)), zap.Any("resonse", toMap(resp)))

				switch msg := resp.MessageResponse.(type) {
				case *pbreflect.ServerReflectionResponse_AllExtensionNumbersResponse:
					r := msg.AllExtensionNumbersResponse
					origSymbol := origReq.MessageRequest.(*pbreflect.ServerReflectionRequest_AllExtensionNumbersOfType).AllExtensionNumbersOfType
					zl.Info("all extensions number",
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
						errorCheck(zl, "unmarshal file descriptor proto error", err)

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
					zl.Warn("received reflection error response", zap.Any("response", toMap(msg.ErrorResponse)))
				default:
					zl.Warn("an unpextec response type was received but not handled")
				}
			}
		default:
			errorCheck(zl, "wuut, invalid response to the request we made error", fmt.Errorf("we received type %T %+v", msg, msg))
		}

		errorCheck(zl, "close send error", stream.CloseSend())
	}

	return nil
}

func toMap(any interface{}) map[string]interface{} {
	cnt, err := json.Marshal(any)
	errorCheck(zlog, "marshal response", err)

	out := map[string]interface{}{}
	errorCheck(zlog, "unmarshal response", json.Unmarshal(cnt, &out))

	return out
}
