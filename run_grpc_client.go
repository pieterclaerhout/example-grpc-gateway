package main

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pieterclaerhout/example-grpc-gateway/example"
	"github.com/pieterclaerhout/go-log"
	"google.golang.org/grpc"
)

func runGRPCClient(serverAddress string) error {

	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := runGRPCClientForYourService(conn); err != nil {
		return err
	}

	if err := runGRPCClientForAnotherService(conn); err != nil {
		return err
	}

	return nil

}

func runGRPCClientForYourService(conn *grpc.ClientConn) error {

	client := example.NewYourServiceClient(conn)

	result, err := client.Echo(context.Background(), &example.StringMessage{
		Value: "Hello From Your Service",
	})
	if err != nil {
		return err
	}

	log.Info("your-service:", result.GetValue())

	return nil

}

func runGRPCClientForAnotherService(conn *grpc.ClientConn) error {

	client := example.NewAnotherServiceClient(conn)

	result, err := client.HelloWorld(context.Background(), &empty.Empty{})
	if err != nil {
		return err
	}

	log.Info("another-service:", result.GetValue())

	return nil

}
