package main

import (
	"flag"
	"fmt"
	"log"

	proto "example/challenges/internal/proto"
)

func main() {
	flag.Parse()
	port := flag.Int("port", 5001, "The server port")

	s, lis, err := proto.NewTasksGrpcServer(fmt.Sprintf("localhost:%d", *port))

	if err != nil {
		log.Fatalf("failed to build server: %v", err)
	}

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
