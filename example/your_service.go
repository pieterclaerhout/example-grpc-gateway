package example

import (
	context "context"
	"net"

	"github.com/pieterclaerhout/go-log"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

// YourService ...
type YourService struct {
}

// New creates new server greeter
func New() *YourService {
	return &YourService{}
}

// Start starts server
func (s *YourService) Start() error {
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	RegisterYourServiceServer(grpcServer, s)
	grpcServer.Serve(lis)
	return nil
}

// Echo just sends back the input
func (s *YourService) Echo(ctx context.Context, msg *StringMessage) (*StringMessage, error) {
	p, _ := peer.FromContext(ctx)
	log.InfoDump(p, "p")
	log.InfoDump(msg, "msg")
	return msg, nil
}
