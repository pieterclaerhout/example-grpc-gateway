package main

import (
	"github.com/lucperkins/rek"
	"github.com/pieterclaerhout/example-grpc-gateway/example"
	"github.com/pieterclaerhout/go-log"
)

func runRestClient(serverAddress string) error {

	if err := runRestClientForYourService(serverAddress); err != nil {
		return err
	}

	if err := runRestClientForAnotherService(serverAddress); err != nil {
		return err
	}

	return nil

}

func runRestClientForYourService(serverAddress string) error {

	res, err := rek.Post(
		"http://localhost:8080/v1/example/echo",
		rek.Json(&example.StringMessage{
			Value: "Hello From Your Service",
		}),
	)
	if err != nil {
		return err
	}

	body, _ := rek.BodyAsString(res.Body())

	log.Info("your-service:", body)

	return nil

}

func runRestClientForAnotherService(serverAddress string) error {

	res, err := rek.Get(
		"http://localhost:8080/v1/another/hello",
	)
	if err != nil {
		return err
	}

	body, _ := rek.BodyAsString(res.Body())

	log.Info("another-service:", body)

	return nil

}
