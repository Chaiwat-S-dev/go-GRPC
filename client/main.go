package main

import (
	"log"

	"github.com/Chaiwat-S-dev/go-gRPC/client/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	creds := insecure.NewCredentials()
	cc, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}
	calculatorClient := services.NewCalculatorClient(cc)
	calculatorService := services.NewCalculatorService(calculatorClient)

	err = calculatorService.Hello("Bond")
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			log.Printf("[%v] %v", grpcErr.Code(), grpcErr.Message())
		}
		log.Fatal(err)
	}
	err = calculatorService.Fibonacci(6)
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			log.Printf("[%v] %v", grpcErr.Code(), grpcErr.Message())
		}
		log.Fatal(err)
	}
	err = calculatorService.Average(1, 2, 3, 4, 5)
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			log.Printf("[%v] %v", grpcErr.Code(), grpcErr.Message())
		}
		log.Fatal(err)
	}
	err = calculatorService.Sum(1, 2, 3, 4, 5)
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			log.Printf("[%v] %v", grpcErr.Code(), grpcErr.Message())
		}
		log.Fatal(err)
	}
}
