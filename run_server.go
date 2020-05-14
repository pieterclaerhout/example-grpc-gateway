package main

import (
	"context"
	"fmt"
	"mime"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pieterclaerhout/example-grpc-gateway/example"
	"github.com/pieterclaerhout/go-log"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
)

func getOpenAPIHandler() http.Handler {
	mime.AddExtensionType(".svg", "image/svg+xml")

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal("creating OpenAPI filesystem:", err.Error())
	}

	return http.FileServer(statikFS)
}

func runServer() error {

	lis, err := net.Listen("tcp", *serverEndpoint)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	example.RegisterYourServiceServer(s, example.NewYourService())
	example.RegisterAnotherServiceServer(s, example.NewAnotherService())

	go func() {
		log.Info("Starting gRPC server:", *serverEndpoint)
		log.Fatal(s.Serve(lis))
	}()

	dialAddr := fmt.Sprintf("dns:///%s", *serverEndpoint)

	conn, err := grpc.DialContext(
		context.Background(),
		dialAddr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}

	jsonpb := &runtime.JSONPb{
		EmitDefaults: true,
		Indent:       "  ",
		OrigName:     true,
	}
	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, jsonpb),
	)

	if err := example.RegisterYourServiceHandler(context.Background(), gwmux, conn); err != nil {
		return err
	}
	if err := example.RegisterAnotherServiceHandler(context.Background(), gwmux, conn); err != nil {
		return err
	}

	oa := getOpenAPIHandler()

	gatewayAddr := "0.0.0.0:8080"

	gwServer := &http.Server{
		Addr: gatewayAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info(r.URL.Path)
			if strings.HasPrefix(r.URL.Path, "/v1") {
				gwmux.ServeHTTP(w, r)
				return
			}
			oa.ServeHTTP(w, r)
		}),
	}

	log.Info("Starting gRPC gateway:", gatewayAddr)
	return gwServer.ListenAndServe()

}
