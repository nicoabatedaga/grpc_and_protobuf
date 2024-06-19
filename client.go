package main

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"

	pb "github.com/nicoabatedaga/grpc_and_protobuf/protos/todo"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewTodoServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Add a task
	task := &pb.Task{Id: "1", Description: "Learn gRPC"}
	res, err := c.AddTask(ctx, task)
	if err != nil {
		log.Fatalf("Could not add task: %v", err)
	}
	log.Printf("AddTask: %s", res.Message)

	// List tasks
	tasks, err := c.ListTasks(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("Could not list tasks: %v", err)
	}
	for _, t := range tasks.Tasks {
		log.Printf("Task: %s, Description: %s", t.Id, t.Description)
	}

	// Delete a task
	res, err = c.DeleteTask(ctx, &pb.TaskId{Id: "1"})
	if err != nil {
		log.Fatalf("Could not delete task: %v", err)
	}
	log.Printf("DeleteTask: %s", res.Message)
}
