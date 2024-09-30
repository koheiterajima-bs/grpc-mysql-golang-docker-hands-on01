package main

import (
	"context"
	"database/sql"
	"log"
	"net"

	_ "github.com/go-sql-driver/mysql"
	pb "github.com/koheiterajima-bs/grpc-mysql-golang-docker-hands-on01/proto/pb"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTodoServiceServer
	db *sql.DB
}

func (s *server) CreateTodo(ctx context.Context, todo *pb.Todo) (*pb.TodoResponse, error) {
	_, err := s.db.Exec("INSERT INTO todos (title, description) VALUES (?, ?)", todo.Title, todo.Description)
	if err != nil {
		return nil, err
	}
	return &pb.TodoResponse{Message: "Todo created"}, nil
}

func (s *server) GetTodo(ctx context.Context, req *pb.TodoRequest) (*pb.Todo, error) {
	row := s.db.QueryRow("SELECT id, description FROM todos WHERE id = ?", req.Id)
	var todo pb.Todo
	err := row.Scan(&todo.Id, &todo.Title, &todo.Description)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (s *server) ListTodos(ctx context.Context, empty *pb.Empty) (*pb.Todos, error) {
	rows, err := s.db.Query("SELECT id, title, description FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos pb.Todos
	for rows.Next() {
		var todo pb.Todo
		if err := rows.Scan(&todo.Id, &todo.Title, &todo.Description); err != nil {
			return nil, err
		}
		todos.Todos = append(todos.Todos, &todo)
	}
	return &todos, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db, err := sql.Open("mysql", "geouser:geouserpassword@tcp(db:3306)/todo_db")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTodoServiceServer(grpcServer, &server{db: db})

	log.Println("gRPC server is running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
