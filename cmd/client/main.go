package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"grpc-greeter/pkg/proto/generated"

	"google.golang.org/grpc"
)

func main() {
	// Connect to gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	client := generated.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// Make request
	resp, err := client.SayHello(ctx, &generated.HelloRequest{Name: "Zee"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	fmt.Printf("Greeting: %s\n", resp.GetMessage())
}
