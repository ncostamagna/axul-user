# axul-user

# Create Project (gRPC on standby)

Command to create grpc file

```sh
protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative  pkg/grpc/userpb/user.proto
```
