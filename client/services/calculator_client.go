package services

import (
	"context"
	"fmt"
	"io"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type CalculatorService interface {
	Hello(name string) error
	Fibonacci(n uint32) error
	Average(numbers ...float64) error
	Sum(numbers ...int32) error
}

type calculatorService struct {
	calculatorClient CalculatorClient
}

func NewCalculatorService(calculatorClient CalculatorClient) CalculatorService {
	return calculatorService{calculatorClient}
}

func (c calculatorService) Hello(name string) error {
	fmt.Println("[client]: Start Hello")
	req := HelloRequest{
		Name:        name,
		CreatedDate: timestamppb.Now(),
	}
	res, err := c.calculatorClient.Hello(context.Background(), &req)
	if err != nil {
		return err
	}
	fmt.Printf("[Client] Request: %v\n", req.Name)
	fmt.Printf("[Client] Response: %v\n", res.Result)
	return nil
}

func (c calculatorService) Fibonacci(n uint32) error {
	fmt.Println("[client]: Start Fibonacci")
	req := FibonacciRequest{
		N: n,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	stream, err := c.calculatorClient.Fibonacci(ctx, &req)
	if err != nil {
		return err
	}
	fmt.Printf("[client] Request: %v\n", req.N)
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Printf("[client] Response: %v\n", res.Result)
	}
	return nil
}

func (c calculatorService) Average(numbers ...float64) error {
	fmt.Println("[client]: Start Average")
	stram, err := c.calculatorClient.Average(context.Background())
	if err != nil {
		return err
	}
	for _, number := range numbers {
		req := AverageRequest{
			Number: number,
		}
		fmt.Printf("[client] Request: %v\n", req.Number)
		stram.Send(&req)
		time.Sleep(time.Second)
	}
	res, err := stram.CloseAndRecv()
	if err != nil {
		return err
	}
	fmt.Printf("[client] Response: %v\n", res.Result)
	return nil
}

func (c calculatorService) Sum(numbers ...int32) error {
	fmt.Println("[client]: Start Sum")
	stream, err := c.calculatorClient.Sum(context.Background())
	if err != nil {
		return err
	}
	go func() {
		for _, number := range numbers {
			req := SumRequest{
				Number: number,
			}
			fmt.Printf("[client] Request: %v\n", req.Number)
			stream.Send(&req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()
	done := make(chan bool)
	errs := make(chan error)
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				errs <- err
			}
			fmt.Printf("[client] Response: %v\n", res.Result)
		}
		done <- true
	}()
	select {
	case <-done:
		return nil
	case err := <-errs:
		return err
	}
}
