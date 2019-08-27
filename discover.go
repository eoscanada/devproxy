package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/eoscanada/devproxy/insecure"
	"github.com/gogo/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pbreflect "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

func dialOptions(endpoint string) (opts []grpc.DialOption) {
	if strings.Contains(endpoint, "*") {
		zlog.Info("with transport credentials")
		creds := credentials.NewTLS(&tls.Config{
			RootCAs:            insecure.CertPool,
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
		errorCheck("dialing to service "+srv, err)

		client := pbreflect.NewServerReflectionClient(conn)
		stream, err := client.ServerReflectionInfo(context.Background())
		errorCheck("setting up client", err)

		err = stream.Send(&pbreflect.ServerReflectionRequest{
			Host:           target,
			MessageRequest: &pbreflect.ServerReflectionRequest_ListServices{ListServices: "*"},
		})
		errorCheck("sending reflection request", err)

		resp, err := stream.Recv()
		errorCheck("receiving reflection response", err)

		cnt, err := json.MarshalIndent(resp, "", "  ")
		errorCheck("json marshal", err)

		fmt.Println("RESP", string(cnt))

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
				errorCheck("sending reflection request", err)

				err = stream.Send(&pbreflect.ServerReflectionRequest{
					Host:           target,
					MessageRequest: &pbreflect.ServerReflectionRequest_AllExtensionNumbersOfType{AllExtensionNumbersOfType: serviceName},
				})
				errorCheck("sending reflection request", err)

				reqs += 2
			}

			//
			for i := 0; i < reqs; i++ {
				resp, err := stream.Recv()
				errorCheck("receiving reflection response", err)

				cnt, err := json.MarshalIndent(resp, "", "  ")
				errorCheck("marshal response", err)

				origReq := resp.OriginalRequest
				switch msg := resp.MessageResponse.(type) {

				case *pbreflect.ServerReflectionResponse_AllExtensionNumbersResponse:
					r := msg.AllExtensionNumbersResponse
					origSymbol := origReq.MessageRequest.(*pbreflect.ServerReflectionRequest_AllExtensionNumbersOfType).AllExtensionNumbersOfType
					fmt.Println("All extensions number", r.BaseTypeName, r.ExtensionNumber, origSymbol)
					conf.extensionNumbers[origSymbol] = resp

				case *pbreflect.ServerReflectionResponse_FileDescriptorResponse:
					r := msg.FileDescriptorResponse

					var filenames []string
					for _, descFile := range r.FileDescriptorProto {
						desc := &descriptor.FileDescriptorProto{}
						err = proto.Unmarshal(descFile, desc)
						errorCheck("unmarshal file descriptor proto", err)

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
						//fmt.Println("File descriptor resp", len(r.FileDescriptorProto), origSymbol)
						conf.fileContainingSymbol[origSymbol] = resp

					case *pbreflect.ServerReflectionRequest_FileByFilename:
						origFilename := origPayload.FileByFilename
						conf.filesByFilename[origFilename] = resp
					}

				case *pbreflect.ServerReflectionResponse_ErrorResponse:
					fmt.Println("ERROR DUDE", msg.ErrorResponse.ErrorMessage)
				default:
					fmt.Println("Some other weird response!")
				}
				_ = cnt
				//fmt.Println("MAMAMMM", string(cnt))

			}
		default:
			errorCheck("wuut, invalid response to the request we made", fmt.Errorf("we received type %T %+v", msg, msg))
		}

		errorCheck("close send", stream.CloseSend())
	}

	return nil
}
