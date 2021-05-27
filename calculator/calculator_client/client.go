package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/LibenHailu/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Calcuator client  ")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not found connect %v ", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	// fmt.Printf("created client : %f", c)
	// doUnary(c)
	// doServerStream(c)
	doClientStream(c)

}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("starting to do a sum unary RPC...")
	req := &calculatorpb.SumRequest{
		FirstNumber:  5,
		SecondNumber: 40,
	}
	res, err := c.Sum(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}

	log.Printf("Response form Sum: %v", res.SumResult)
}

func doServerStream(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("starting to do a primeDecomposition server streaming RPC...")
	req := &calculatorpb.PrimaryNumberDecompositionRequest{
		Number: 12,
	}
	stream, err := c.PrimeNumberDecomposition(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happend: %v", err)
		}
		fmt.Println(res.GetPrimeFactor())
	}

}

func doClientStream(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("starting to do a compute average client streaming RPC...")
	stream, err := c.ComputeAverage(context.Background())

	if err != nil {
		log.Fatalf("Error while opening stream: %v", err)
	}

	numbers := []int32{3, 5, 9, 54, 23}

	for _, number := range numbers {
		fmt.Printf("Sending number :%v\n", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while reciveing reponse: %v", err)
	}

	fmt.Printf("The average is: %v\n ", res.GetAverage())
}
