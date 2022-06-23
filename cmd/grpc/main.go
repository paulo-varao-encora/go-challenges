package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "example/challenges/cmd/grpc/tasks"
	"example/challenges/internal"
	"example/challenges/internal/mux"
)

type server struct {
	pb.UnimplementedTasksGrpcServer
	table internal.TaskTable
}

func (s *server) RetrieveAll(ctx context.Context, in *pb.Empty) (*pb.TaskList, error) {
	tasks, err := s.table.RetrieveAll()

	if err != nil {
		return nil, status.Newf(codes.Internal, err.Error()).Err()
	}

	return getProtoTaskList(tasks), nil
}

func (s *server) FilterTasks(ctx context.Context, in *pb.FilterRequest) (*pb.TaskList, error) {
	completed := in.GetCompleted()
	tasks, err := s.table.Filter(completed)

	if err != nil {
		return nil, status.Newf(codes.Internal, err.Error()).Err()
	}

	return getProtoTaskList(tasks), nil
}

func main() {
	flag.Parse()
	port := flag.Int("port", 5001, "The server port")
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv, err := newTableServer()
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTasksGrpcServer(s, srv)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func getProtoTaskList(tasks []internal.Task) *pb.TaskList {
	result := []*pb.ExistingTask{}

	for _, t := range tasks {
		result = append(result, &pb.ExistingTask{ID: t.ID, Name: t.Name, Completed: t.Completed})
	}

	return &pb.TaskList{Tasks: result}
}

func newTableServer() (*server, error) {
	table, err := mux.SelectDBImpl()

	if err != nil {
		return nil, err
	}
	srv := server{table: table}
	return &srv, nil
}
