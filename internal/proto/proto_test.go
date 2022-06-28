package proto

import (
	"context"
	"net"
	"reflect"
	"testing"
	"time"

	"example/challenges/internal"
	pb "example/challenges/internal/proto/tasks"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGrpcServer(t *testing.T) {

	addr := "localhost:4990"

	// Server
	s, lis := newServer(t, addr)
	defer s.Stop()
	defer lis.Close()

	// Client
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Errorf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTasksGrpcClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	t.Run("list all tasks", func(t *testing.T) {
		r, err := c.RetrieveAll(ctx, &pb.Empty{})
		if err != nil {
			t.Errorf("failed to retrieve all tasks: %v", err)
		}

		got := getInternalTaskList(r.GetTasks())
		assertTasks(t, got, internal.DefaultTasks)
	})

	t.Run("create a new task and returns its ID", func(t *testing.T) {
		newTask := pb.NewTask{Name: "Sample task", Completed: true}

		r, err := c.Create(ctx, &newTask)
		if err != nil {
			t.Errorf("failed to create task: %v", err)
		}

		if r.GetId() < 5 {
			t.Errorf("got %v want greater than 4", r.GetId())
		}
	})

	t.Run("get single task by its ID", func(t *testing.T) {
		id := int64(1)
		taskID := pb.TaskID{Id: id}

		r, err := c.RetrieveTaskByID(ctx, &taskID)

		if err != nil {
			t.Errorf("failed to retrieve task %v: %v", id, err)
		}

		got := internal.Task{ID: r.GetID(), Name: r.GetName(), Completed: r.GetCompleted()}
		want := internal.DefaultTasks[0]

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("update single task by its id", func(t *testing.T) {
		id := int64(5)
		existingTask := pb.ExistingTask{ID: id, Name: "Sample task updated", Completed: false}
		_, err := c.Update(ctx, &existingTask)

		if err != nil {
			t.Errorf("failed to update task %v: %v", id, err)
		}
	})

	t.Run("list all completed tasks", func(t *testing.T) {
		req := pb.FilterRequest{Completed: true}

		r, err := c.FilterTasks(ctx, &req)

		if err != nil {
			t.Errorf("failed to filter tasks: %v", err)
		}

		got := getInternalTaskList(r.GetTasks())
		gotSize := len(got)

		if gotSize != 2 {
			t.Errorf("got %v, expected 2", gotSize)
		}
	})
}

func assertTasks(t testing.TB, got, want []internal.Task) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func getInternalTaskList(protoTasks []*pb.ExistingTask) []internal.Task {
	result := []internal.Task{}

	for _, t := range protoTasks {
		result = append(result, internal.Task{ID: t.ID, Name: t.Name, Completed: t.Completed})
	}

	return result
}

func newServer(t testing.TB, address string) (*grpc.Server, net.Listener) {
	t.Helper()
	s, lis, err := NewTasksGrpcServer(address)

	if err != nil {
		t.Errorf("failed to build server: %v", err)
	}

	go s.Serve(lis)

	return s, lis
}
