package main

import (
	"context"
	"log"
	"time"

	pb "github.com/koheiterajima-bs/grpc-mysql-golang-docker-hands-on01/proto/pb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTodoServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = client.CreateTodo(ctx, &pb.Todo{title: "Learn gRPC", Description: "gRPC with Go and MySQL"})
	if err != nil {
		log.Fatalf("could not create todo: %v", err)
	}

	log.Println("Todo created successfully")
}
