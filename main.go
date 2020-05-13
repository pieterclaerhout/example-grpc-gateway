package main

import (
	"context"
	"flag"
	"fmt"
	"mime"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/lucperkins/rek"
	"github.com/pieterclaerhout/example-grpc-gateway/example"
	"github.com/pieterclaerhout/go-log"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"

	_ "github.com/pieterclaerhout/example-grpc-gateway/statik"
)

func getOpenAPIHandler() http.Handler {
	mime.AddExtensionType(".svg", "image/svg+xml")

	statikFS, err := fs.New()
	if err != nil {
		panic("creating OpenAPI filesystem: " + err.Error())
	}

	return http.FileServer(statikFS)
}

var serverEndpoint = flag.String("endpoint", "0.0.0.0:9090", "gRPC server endpoint")
var what = flag.String("what", "server", "server|client")

func runServer() error {

	log.Info("Creating gRPC listener")
	lis, err := net.Listen("tcp", *serverEndpoint)
	if err != nil {
		return err
	}

	log.Info("Creating gRPC server")
	s := grpc.NewServer()
	example.RegisterYourServiceServer(s, example.New())

	go func() {
		log.Info("Starting gRPC server")
		log.Fatal(s.Serve(lis))
	}()

	dialAddr := fmt.Sprintf("dns:///%s", *serverEndpoint)

	log.Info("Creating gRPC connection", dialAddr)
	conn, err := grpc.DialContext(
		context.Background(),
		dialAddr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}

	log.Info("Creating mux")
	jsonpb := &runtime.JSONPb{
		EmitDefaults: true,
		Indent:       "  ",
		OrigName:     true,
	}
	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, jsonpb),
	)
	err = example.RegisterYourServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		return err
	}

	oa := getOpenAPIHandler()

	log.Info("Creating gateway server")

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

	log.Info("Serving gRPC-Gateway and OpenAPI Documentation on http://", gatewayAddr)
	return gwServer.ListenAndServe()

}

func runGRPCClient() error {

	conn, err := grpc.Dial(*serverEndpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := example.NewYourServiceClient(conn)

	result, err := client.Echo(context.Background(), &example.StringMessage{
		Value: "Hello world",
	})
	if err != nil {
		return err
	}

	log.InfoDump(result, "result")

	return nil

}

func runRestClient() error {

	res, err := rek.Post(
		"http://localhost:8080/v1/example/echo",
		rek.Json(&example.StringMessage{
			Value: "Hello world",
		}),
		rek.Timeout(5*time.Second),
	)
	if err != nil {
		return err
	}

	log.Info(rek.BodyAsString(res.Body()))

	return nil
}

func main() {

	log.PrintColors = true
	log.PrintTimestamp = true

	flag.Parse()

	log.Info("Running:", *what)

	var err error

	switch *what {
	case "server":
		err = runServer()
	case "grpc-client":
		err = runGRPCClient()
	case "rest-client":
		err = runRestClient()
	default:
		log.Fatal("Incorrect option for argument what:", *what)
	}

	log.CheckError(err)

}
