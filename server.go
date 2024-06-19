package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "github.com/nicoabatedaga/grpc_and_protobuf/protos/todo"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTodoServiceServer
	mu    sync.Mutex
	tasks map[string]*pb.Task
}

func (s *server) AddTask(ctx context.Context, task *pb.Task) (*pb.TaskResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tasks[task.Id] = task
	return &pb.TaskResponse{Message: "Task added"}, nil
}

func (s *server) ListTasks(ctx context.Context, empty *pb.Empty) (*pb.Tasks, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tasks := make([]*pb.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return &pb.Tasks{Tasks: tasks}, nil
}

func (s *server) DeleteTask(ctx context.Context, taskId *pb.TaskId) (*pb.TaskResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.tasks, taskId.Id)
	return &pb.TaskResponse{Message: "Task deleted"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTodoServiceServer(s, &server{tasks: make(map[string]*pb.Task)})

	log.Printf("Server is running on port :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
