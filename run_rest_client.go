package main

import (
	"github.com/lucperkins/rek"
	"github.com/pieterclaerhout/example-grpc-gateway/example"
	"github.com/pieterclaerhout/go-log"
)

func runRestClient() error {

	if err := runRestClientForYourService(); err != nil {
		return err
	}

	if err := runRestClientForAnotherService(); err != nil {
		return err
	}

	return nil

}

func runRestClientForYourService() error {

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

func runRestClientForAnotherService() error {

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
