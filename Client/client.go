package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"calculator/calcpb"
)

func main() {
	cc, err := grpc.Dial("localhost:1010", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.close()

	c := calcpb.NewCalculatorServiceClient(cc)

	// Sum function
	// Sum(c)

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
		Sum: &calcpb.CalculatorRequest{
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
		Number: &calcpb.Number{
			p: 20,
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

	requests := []*calcpb.ComputeAverageRequest{}{
		&calcpb.ComputeAverageRequest{
			Number: &calcpb.Number{
				P: 4,
			},
		},
		&calcpb.ComputeAverageRequest{
			Number: &calcpb.Number{
				P: 5,
			},
		},
		&calcpb.ComputeAverageRequest{
			Number: &calcpb.Number{
				P: 9,
			},
		},
		&calcpb.ComputeAverageRequest{
			Number: &calcpb.Number{
				P: 13,
			},
		},
		&calcpb.ComputeAverageRequest{
			Number: &calcpb.Number{
				P: 25,
			},
		},
		&calcpb.ComputeAverageRequest{
			Number: &calcpb.Number{
				P: 4,
			},
		},
	}

	for _, req := range requests{
		fmt.Println("\nSending Request...", req)
		stream.Send(req)
		time.Sleen(1*time.Second)
	}

	resp, err := stream.CloseAndRecv()
	if err!=nil{
		log.Fatalf("Error while receiving response from server: %v", err)
	}
	fmt.Println("\n Computed Average: ", resp.GetResult())
 }

 func FindMaxNumber(c calcpb.CalculatorServiceClient){
	fmt.Println("Finding the Max number over the incoming stream of numbers...")

	requests := []*calcpb.FindMaxNumberRequest{
		&calcpb.FindMaxNumberRequest{
			Number: &calcpb.Number{
				P: 1,
			},
		},
		&calcpb.FindMaxNumberRequest{
			Number: &calcpb.Number{
				P: 3,
			},
		},
		&calcpb.FindMaxNumberRequest{
			Number: &calcpb.Number{
				P: 7,
			},
		},
		&calcpb.FindMaxNumberRequest{
			Number: &calcpb.Number{
				P: 9,
			},
		},
		&calcpb.FindMaxNumberRequest{
			Number: &calcpb.Number{
				P: 2,
			},
		},
		&calcpb.FindMaxNumberRequest{
			Number: &calcpb.Number{
				P: 5,
			},
		},
		&calcpb.FindMaxNumberRequest{
			Number: &calcpb.Number{
				P: 22,
			},
		},
		&calcpb.FindMaxNumberRequest{
			Number: &calcpb.Number{
				P: 15,
			},
		},
		&calcpb.FindMaxNumberRequest{
			Number: &calcpb.Number{
				P: 21,
			},
		},
		&calcpb.FindMaxNumberRequest{
			Number: &calcpb.Number{
				P: 19,
			},
		},
	}

	stream, err := c.FindMaxNumber(context.Background())
	if err!=nil{
		log.Fatalf("error occured while performing client side streaming: %v", err)
	}

	waitchan := make(chan struct{})

	go func(requests []*calcpb.FindMaxNumberRequest){
		for _, req := range requests{
			fmt.Println("\n Sending Request... : ", req.Calculator)
			err := stream.Send(req)
			if err!=nil{
				log.Fatalf("error while sending request to FindMaxNUmber service: %v", err)
			}
			time.Sleep(1000*time.Millisecond)
		}(requests)

		go func(){
			for {
				resp, err := stream.Recv()
				if err==io.EOF{
					close(waitchan)
					return
				}
				if err !=nil{
					log.Fatalf("error while receiving response from server: %v", err)
				}
				fmt.Printf("\n Max Element till now in the Stream: %v", resp.GetResult())
			}
		}()

		<-waitchan
	}
 }
