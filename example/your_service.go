package example

import (
	context "context"

	empty "github.com/golang/protobuf/ptypes/empty"
)

type YourService struct {
}

func NewYourService() *YourService {
	return &YourService{}
}

func (s *YourService) Echo(c context.Context, msg *StringMessage) (*StringMessage, error) {
	return msg, nil
}

type AnotherService struct {
}

func NewAnotherService() *AnotherService {
	return &AnotherService{}
}

func (s *AnotherService) HelloWorld(c context.Context, empty *empty.Empty) (*StringMessage, error) {
	return &StringMessage{Value: "Hello World"}, nil
}
