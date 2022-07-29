package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"time"

	"calculator/calcpb"

	"google.golang.org/grpc"
)

type server struct {
	calcpb.UnimplementedCalculatorServiceServer
}

func (*server) Sum(ctx context.Context, req *calcpb.SumRequest) (resp *calcpb.SumResponse, err error) {
	fmt.Println("Implementing Sum function of Calculator")

	a := req.GetS().GetA()
	b := req.GetS().GetB()

	result := a + b

	res := result

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

func (*server) PrimeNumber(req *calcpb.PrimeNumberRequest, resp calcpb.CalculatorService_PrimeNumberServer) error {
	fmt.Println("Implementing Prime Numbers Operation")

	p := req.GetP().GetP()

	for i := 0; int64(i) < p; i++ {
		isp := isprime(i)
		if isp {
			result := int64(i)
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
	count := int64(0)
	res := int64(0)

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			result = float64(res / count)
			return stream.SendAndClose(&calcpb.ComputeAverageResponse{
				Result: result,
			})
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		count += 1
		res += msg.GetP().GetP()
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
		max := int64(0)
		num := req.GetP().GetP()

		if max < num {
			result := max
			sendErr := stream.Send(&calcpb.FindmaxNumberResponse{
				Result: result,
			})
			if sendErr != nil {
				log.Fatalf("error while sending response to FindMacNumber Client: %v", err)
				return err
			}
		}
	}
}

func main() {
	fmt.Println("Server is Listening...")

	listen, err := net.Listen("tcp", "0.0.0.0:1010")
	if err != nil {
		log.Fatal("Failed to Listen: %v", err)
	}

	s := grpc.NewServer()
	calcpb.RegisterCalculatorServiceServer(s, &server{})

	if err = s.Serve(listen); err != nil {
		log.Fatal("Failed to Serve: %v", err)
	}
}
