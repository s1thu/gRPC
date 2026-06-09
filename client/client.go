package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	pb "github.com/s1thu/gRPC/todo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const defaultServerAddress = "localhost:50051"

func main() {
	serverAddress := flag.String("server", defaultServerAddress, "gRPC server address")
	action := flag.String("action", "list", "RPC to call: create, delete, modify, list")
	name := flag.String("name", "", "todo name")
	description := flag.String("description", "", "todo description")
	done := flag.Bool("done", false, "todo completion state")
	id := flag.String("id", "", "todo id")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, *serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewTodoServiceClient(conn)

	switch strings.ToLower(strings.TrimSpace(*action)) {
	case "create":
		if strings.TrimSpace(*name) == "" {
			log.Fatal("-name is required for create")
		}
		createTodo(ctx, client, *name, *description, *done)
	case "delete":
		if strings.TrimSpace(*id) == "" {
			log.Fatal("-id is required for delete")
		}
		deleteTodo(ctx, client, *id)
	case "modify":
		if strings.TrimSpace(*id) == "" {
			log.Fatal("-id is required for modify")
		}
		modifyTodo(ctx, client, *id, *name, *description, *done)
	case "list":
		listTodos(ctx, client)
	default:
		log.Fatalf("unknown action %q", *action)
	}
}

func createTodo(ctx context.Context, client pb.TodoServiceClient, name, description string, done bool) {
	resp, err := client.CreateTodo(ctx, &pb.NewTodo{
		Name:        name,
		Description: description,
		Done:        done,
	})
	if err != nil {
		log.Fatalf("create todo failed: %v", err)
	}
	printTodo("created", resp)
}

func deleteTodo(ctx context.Context, client pb.TodoServiceClient, todoID string) {
	resp, err := client.DeleteTodo(ctx, &pb.TodoId{Id: todoID})
	if err != nil {
		log.Fatalf("delete todo failed: %v", err)
	}
	fmt.Printf("deleted todo %s: %v\n", todoID, resp)
}

func modifyTodo(ctx context.Context, client pb.TodoServiceClient, todoID, name, description string, done bool) {
	resp, err := client.ModifyTodo(ctx, &pb.Todo{
		Id:          todoID,
		Name:        name,
		Description: description,
		Done:        done,
	})
	if err != nil {
		log.Fatalf("modify todo failed: %v", err)
	}
	printTodo("updated", resp)
}

func listTodos(ctx context.Context, client pb.TodoServiceClient) {
	stream, err := client.ListTodos(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("list todos failed: %v", err)
	}

	count := 0
	for {
		todo, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("receive todo failed: %v", err)
		}
		count++
		printTodo(fmt.Sprintf("todo %d", count), todo)
	}

	if count == 0 {
		fmt.Fprintln(os.Stdout, "no todos found")
	}
}

func printTodo(label string, todo *pb.Todo) {
	if todo == nil {
		fmt.Printf("%s: <nil>\n", label)
		return
	}
	fmt.Printf("%s: id=%s name=%q description=%q done=%t\n", label, todo.GetId(), todo.GetName(), todo.GetDescription(), todo.GetDone())
}
