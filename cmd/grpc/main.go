package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"example/challenges/internal"
	"example/challenges/internal/mux"
	pb "example/challenges/internal/proto/tasks"
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

func (s *server) RetrieveTaskByID(ctx context.Context, in *pb.TaskID) (*pb.ExistingTask, error) {
	id := in.GetId()
	t, err := s.table.FindByID(id)

	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, status.Newf(codes.InvalidArgument, err.Error()).Err()
		}
		return nil, status.Newf(codes.Internal, err.Error()).Err()
	}

	return &pb.ExistingTask{ID: t.ID, Name: t.Name, Completed: t.Completed}, nil
}

func (s *server) Create(ctx context.Context, in *pb.NewTask) (*pb.TaskID, error) {
	task := internal.Task{Name: in.GetName(), Completed: in.GetCompleted()}
	id, err := s.table.Create(task)

	if err != nil {
		return nil, status.Newf(codes.Internal, err.Error()).Err()
	}

	return &pb.TaskID{Id: id}, nil
}

func (s *server) Update(ctx context.Context, in *pb.ExistingTask) (*pb.Empty, error) {
	task := internal.Task{ID: in.GetID(), Name: in.GetName(), Completed: in.GetCompleted()}
	_, err := s.table.Update(task)

	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, status.Newf(codes.InvalidArgument, err.Error()).Err()
		}
		return nil, status.Newf(codes.Internal, err.Error()).Err()
	}

	return &pb.Empty{}, nil
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
