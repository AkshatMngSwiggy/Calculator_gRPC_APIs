package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"time"

	"calculator/calcpb"
)

type server struct {
	calcpb.UnimplementedGreetServiceServer
}

func (*server) Sum(ctx context.Context, req *calcpb.SumRequest) (resp *calcpb.SumResponse, err error) {
	fmt.Println("Implementing Sum function of Calculator")

	a := req.GetSum().GetA()
	b := req.GetSUm().GetB()

	result := a + b

	res := fmt.Sprintf("Sum of %v and %v equal to %v", a, b, result)

	resp = &calcpb.SumResponse{
		Result: res,
	}

	return resp, nil
}

func isprime(n int) bool {
	if n <= 1 {
		return false
	} else {
		for i := 2; i < int(math.Sqrt(float64(n))+1); i++ {
			if n%i == 0 {
				return false
			}
		}
	}
	return true
}

func (*server) PrimeNumber(ctx context.Context, req *calcpb.PrimeNumberRequest, resp *calcpb.PrimeNumberResponse) error {
	fmt.Println("Implementing Prime Numbers Operation")

	p := req.GetNumber().GetP()

	for i := 0; i < p; i++ {
		isp := isprime(i)
		if isp {
			result := fmt.Sprintf("Prime Number between 0 and %v: %v", p, i)
			res := calcpb.PrimeNumberResponse{
				Result: result,
			}
			time.Sleep(100 * time.Millisecond)
			resp.Send(&res)
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calcpb.CalculatorService_ComputeAverageServer) error {
	fmt.Println("Implementing Computer Average Function")
	result := 0.0
	count := 0.0
	res := 0.0

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			result = res / count
			return stream.SendAndClose(&calcpb.ComputeAverage{
				Result: result,
			})
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		count += 1
		res += msg.GetNumber().GetP()
	}

}

func (*server) FindMaxNumber(stream calcpb.CalculatorService_FindMaxNumberServer) error {
	fmt.Println("Finding Max Number in the stream...")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatal("error while receiving data from FindMaxNUmber client: %v", err)
			return err
		}
		max := 0
		num := req.GetNumber().GetP()

		if max < num {
			result := max
			sendErr := stream.Send(&calcpb.FindmaxNumberResponse{
				Result: result,
			})
		}

		if sendErr != nil {
			log.Fatalf("error while sending response to FindMacNumber Client: %v", err)
			return err
		}
	}
}

func main() {
	fmt.Println("Server is Listening...")

	listen, err = net.listen("tcp", "1.0.1.0:1010")
	if err != nil {
		log.Fatal("Failed to Listen: %v", err)
	}

	s := grpc.NewServer()
	calcpb.RegisterGreetServiceServer(s, &server{})

	if err = s.Serve(listen); err != nil {
		log.Fatal("Failed to Serve: %v", err)
	}
}
