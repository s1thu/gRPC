# gRPC

Learnign gRPC using Go Lang

## Run

Start the server first:

```bash
go run ./server
```

Then run the client from another terminal:

```bash
go run ./client
```

Client usage examples:

```bash
go run ./client -action=list
go run ./client -action=create -name "Buy milk" -description "2 liters"
go run ./client -action=modify -id <todo-id> -name "Updated name" -done=true
go run ./client -action=delete -id <todo-id>
```
