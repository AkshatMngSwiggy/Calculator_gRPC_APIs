package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"calculator/calcpb"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:1010", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := calcpb.NewCalculatorServiceClient(cc)

	// Sum function
	Sum(c)

	// Prime Number
	// PrimeNumber(c)

	// Compute Average
	// ComputeAverage(c)

	// Find Max Number
	// FindMaxNumber(c)
}

func Sum(c calcpb.CalculatorServiceClient) {
	fmt.Println("Starting sending two numbers for Sum")
	req := calcpb.CalculatorRequest{
		Sum: &calcpb.Calculator{
			a: 4,
			b: 5,
		},
	}

	resp, err := c.Calculator(context.Background(), &req)
	if err != nil {
		log.Fatalf("error while calling sum: %v", err)
	}

	log.Printf("Response from Sum Server: %v", resp.Result)
}

func PrimeNumber(c calcpb.CalculatorServiceClient) {
	fmt.Println("Ready to receive prime numbers...")

	req := calcpb.PrimeNumberRequest{
		P: &calcpb.Number{
			P: 20,
		},
	}

	respStream, err := c.PrimeNumber(context.Background(), &req)
	if err != nil {
		log.Fatalf("error while calling PrimeNumber: %v", err)
	}

	for {
		msg, err := respStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while receving server stream: %v", err)
		}
		fmt.Println("Response from PrimeNumber Server: ", msg.Result)
	}
}

func ComputeAverage(c calcpb.CalculatorServiceClient) {
	fmt.Println("Starting Client Side Streaming over computing average...")
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error occured while performing client-side streaming: %v", err)
	}

	requests := []*calcpb.ComputeAverageRequest{
		&calcpb.ComputeAverageRequest{
			P: &calcpb.Number{
				P: 4,
			},
		},
		&calcpb.ComputeAverageRequest{
			P: &calcpb.Number{
				P: 5,
			},
		},
		&calcpb.ComputeAverageRequest{
			P: &calcpb.Number{
				P: 9,
			},
		},
		&calcpb.ComputeAverageRequest{
			P: &calcpb.Number{
				P: 13,
			},
		},
		&calcpb.ComputeAverageRequest{
			P: &calcpb.Number{
				P: 25,
			},
		},
		&calcpb.ComputeAverageRequest{
			P: &calcpb.Number{
				P: 4,
			},
		},
	}

	for _, req := range requests {
		fmt.Println("\nSending Request...", req)
		stream.Send(req)
		time.Sleep(1 * time.Second)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from server: %v", err)
	}
	fmt.Println("\n Computed Average: ", resp.GetResult())
}

func FindMaxNumber(c calcpb.CalculatorServiceClient) {
	fmt.Println("Finding the Max number over the incoming stream of numbers...")

	requests := []*calcpb.FindmaxNumberRequest{
		&calcpb.FindmaxNumberRequest{
			P: &calcpb.Number{
				P: 1,
			},
		},
		&calcpb.FindmaxNumberRequest{
			P: &calcpb.Number{
				P: 3,
			},
		},
		&calcpb.FindmaxNumberRequest{
			P: &calcpb.Number{
				P: 7,
			},
		},
		&calcpb.FindmaxNumberRequest{
			P: &calcpb.Number{
				P: 9,
			},
		},
		&calcpb.FindmaxNumberRequest{
			P: &calcpb.Number{
				P: 2,
			},
		},
		&calcpb.FindmaxNumberRequest{
			P: &calcpb.Number{
				P: 5,
			},
		},
		&calcpb.FindmaxNumberRequest{
			P: &calcpb.Number{
				P: 22,
			},
		},
		&calcpb.FindmaxNumberRequest{
			P: &calcpb.Number{
				P: 15,
			},
		},
		&calcpb.FindmaxNumberRequest{
			P: &calcpb.Number{
				P: 21,
			},
		},
		&calcpb.FindmaxNumberRequest{
			P: &calcpb.Number{
				P: 19,
			},
		},
	}

	stream, err := c.FindMaxNumber(context.Background())
	if err != nil {
		log.Fatalf("error occured while performing client side streaming: %v", err)
	}

	waitchan := make(chan struct{})

	go func(requests []*calcpb.FindmaxNumberRequest) {
		for _, req := range requests {
			fmt.Println("\n Sending Request... : ", req.Calculator)
			err := stream.Send(req)
			if err != nil {
				log.Fatalf("error while sending request to FindMaxNUmber service: %v", err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}(requests)

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				close(waitchan)
				return
			}
			if err != nil {
				log.Fatalf("error while receiving response from server: %v", err)
			}
			fmt.Printf("\n Max Element till now in the Stream: %v", resp.GetResult())
		}
	}()

	<-waitchan
}
