package main

import (
	"context"
	"mime"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pieterclaerhout/example-grpc-gateway/example"
	"github.com/pieterclaerhout/go-log"
	"github.com/rakyll/statik/fs"
	"github.com/soheilhy/cmux"
	"golang.org/x/sync/errgroup"
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

func runServer(serverAddress string) error {

	lis, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		return err
	}

	m := cmux.New(lis)

	grpcListener := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpListener := m.Match(cmux.HTTP1Fast())

	s := grpc.NewServer()
	example.RegisterYourServiceServer(s, example.NewYourService())
	example.RegisterAnotherServiceServer(s, example.NewAnotherService())

	dialAddr := "dns:///0.0.0.0:8080"

	conn, err := grpc.DialContext(
		context.Background(),
		dialAddr,
		grpc.WithInsecure(),
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
			if strings.HasPrefix(r.URL.Path, "/v1") {
				gwmux.ServeHTTP(w, r)
				return
			}
			oa.ServeHTTP(w, r)
		}),
	}

	log.Info("Starting gRPC gateway:", gatewayAddr)

	g := errgroup.Group{}

	g.Go(func() error {
		return s.Serve(grpcListener)
	})

	g.Go(func() error {
		return gwServer.Serve(httpListener)
	})

	g.Go(m.Serve)

	return g.Wait()

}
