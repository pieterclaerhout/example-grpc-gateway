package main

import (
	"flag"

	"github.com/pieterclaerhout/go-log"

	_ "github.com/pieterclaerhout/example-grpc-gateway/statik"
)

var serverEndpoint = flag.String("endpoint", "0.0.0.0:9090", "gRPC server endpoint")
var what = flag.String("what", "server", "server|client")

func main() {

	log.PrintColors = true
	log.PrintTimestamp = true

	flag.Parse()

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
