package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	pb "example/challenges/cmd/grpc/tasks"
	"example/challenges/internal"
	"example/challenges/internal/orm"
	"example/challenges/internal/repository"
)

type server struct {
	pb.UnimplementedTasksGrpcServer
	table internal.TaskTable
}

func (s *server) RetrieveAll(ctx context.Context, in *pb.Empty) (*pb.TaskList, error) {
	tasks, _ := s.table.RetrieveAll()

	result := []*pb.ExistingTask{}

	for _, t := range tasks {
		result = append(result, &pb.ExistingTask{ID: t.ID, Name: t.Name, Completed: t.Completed})
	}

	return &pb.TaskList{Tasks: result}, nil
}

func (s *server) FilterTasks(ctx context.Context, in *pb.FilterRequest) (*pb.TaskList, error) {
	completed := in.GetCompleted()
	tasks, _ := s.table.Filter(completed)

	result := []*pb.ExistingTask{}

	for _, t := range tasks {
		result = append(result, &pb.ExistingTask{ID: t.ID, Name: t.Name, Completed: t.Completed})
	}

	return &pb.TaskList{Tasks: result}, nil
}

func main() {
	flag.Parse()
	port := flag.Int("port", 5001, "The server port")
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	srv := newCrudServer()
	pb.RegisterTasksGrpcServer(s, srv)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func newCrudServer() *server {
	dbImpl := os.Getenv("DB_IMPL")
	var table internal.TaskTable
	var err error

	if dbImpl == "orm" {
		rep, repErr := orm.NewTaskOrm()
		table = &rep
		err = repErr
	} else {
		rep, repErr := repository.NewTaskCrud()
		table = &rep
		err = repErr
	}

	if err != nil {
		log.Fatalf("failed to update table: %v", err)
	}
	srv := server{table: table}
	return &srv
}
