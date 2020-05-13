package example

import (
	context "context"

	"github.com/pieterclaerhout/go-log"
)

// YourService ...
type YourService struct {
}

// New creates new server greeter
func New() *YourService {
	return &YourService{}
}

// Echo just sends back the input
func (s *YourService) Echo(ctx context.Context, msg *StringMessage) (*StringMessage, error) {
	// p, _ := peer.FromContext(ctx)
	log.InfoDump(msg, "msg")
	return msg, nil
}
