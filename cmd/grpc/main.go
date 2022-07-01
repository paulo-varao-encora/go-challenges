package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	proto "example/challenges/internal/proto"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "example/challenges/internal/proto/tasks"
)

func main() {
	flag.Parse()

	// gRPC server on port 5001
	port := flag.Int("port", 5001, "The server port")

	s, lis, err := proto.NewTaskServer(fmt.Sprintf("localhost:%d", *port))

	if err != nil {
		log.Fatalf("failed to build server: %v", err)
	}

	log.Printf("server listening at %v", lis.Addr())
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:5001",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()

	err = pb.RegisterTaskManagerHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	// gRPC gateway server on port 5002
	gwServer := &http.Server{
		Addr:    ":5002",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:5002")
	log.Fatalln(gwServer.ListenAndServe())
}
