package main

import (
	"context"
	"flag"

	"github.com/pieterclaerhout/example-grpc-gateway/example"
	"github.com/pieterclaerhout/go-log"
	"google.golang.org/grpc"
)

var serverEndpoint = flag.String("endpoint", ":8080", "gRPC server endpoint")
var what = flag.String("what", "server", "server|client")

func runServer() error {

	g := example.New()
	err := g.Start()
	return err

	// ctx := context.Background()
	// ctx, cancel := context.WithCancel(ctx)
	// defer cancel()

	// mux := runtime.NewServeMux()
	// opts := []grpc.DialOption{grpc.WithInsecure()}
	// err := example.RegisterYourServiceHandlerFromEndpoint(ctx, mux, *serverEndpoint, opts)
	// if err != nil {
	// 	return err
	// }

	// return http.ListenAndServe(":8081", mux)

}

func runClient() error {

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

func main() {

	flag.Parse()

	log.Info("Running:", *what)

	var err error

	switch *what {
	case "server":
		err = runServer()
	case "client":
		err = runClient()
	default:
		log.Fatal("Incorrect option for argument what:", *what)
	}

	log.CheckError(err)

}
