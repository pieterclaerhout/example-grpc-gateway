package main

import (
	"flag"

	_ "github.com/pieterclaerhout/example-grpc-gateway/statik" // Statik resources
	"github.com/pieterclaerhout/go-log"
)

var serverEndpoint = flag.String("endpoint", "0.0.0.0:8080", "gRPC and HTTP server endpoint")
var what = flag.String("what", "server", "server|client")

func main() {

	log.PrintColors = true
	log.PrintTimestamp = true

	flag.Parse()

	var err error

	switch *what {
	case "server":
		err = runServer(*serverEndpoint)
	case "grpc-client":
		err = runGRPCClient(*serverEndpoint)
	case "rest-client":
		err = runRestClient(*serverEndpoint)
	default:
		log.Fatal("Incorrect option for argument what:", *what)
	}

	log.CheckError(err)

}
