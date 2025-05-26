# grpc-greeter: A Simple Go gRPC Service

## Overview
grpc-greeter is a straightforward example project demonstrating how to build a basic gRPC service in Go. 
It includes a server that implements a "Greeter" service (saying hello) and a client that interacts with it. 
The project is structured to follow common Go project layouts and best practices for gRPC development, 
including Protocol Buffer compilation, dependency management, and basic logging/interceptor patterns.

## Features
- Go gRPC Server: Implements a simple "Greeter" service.
- Protocol Buffers: Defines the service interface and messages using .proto files.
- Code Generation: Automates the generation of Go code from .proto definitions.
- Structured Logging: Uses a custom logger for consistent output.
- gRPC Interceptors: Includes basic logging, recovery, and validation interceptors for server-side middleware.
- Dependency Injection (Basic): Services and use cases are wired together via constructors.
- Configuration Management: Loads server configuration from an external source (assumed config.yaml).
- WSL-Friendly Setup: Designed with common WSL (Windows Subsystem for Linux) development challenges in mind.

## Project Structure
grpc-greeter
├── cmd/
│   ├── client/       # Go gRPC client application
│   │   └── main.go
│   └── server/       # Go gRPC server application
│       └── main.go
├── internal/
│   ├── application/
│   │   └── usecases/ # Application-specific business logic orchestration
│   │       └── greeting_usecase.go
│   ├── config/       # Application configuration loading
│   │   └── config.go
│   ├── domain/
│   │   ├── entities/ # Core business entities
│   │   │   └── greeting.go
│   │   └── services/ # Domain services implementing business rules
│   │       ├── greeting_policies.go
│   │       └── greeting_service.go
│   └── infrastructure/
│       ├── grpc/     # gRPC specific implementations
│       │   ├── handlers/     # gRPC service handlers
│       │   │   └── greeter_handler.go
│       │   └── interceptors/ # gRPC server interceptors/middleware
│       │       └── middleware.go
│       └── logging/  # Custom logging implementation
│           └── logger.go
├── pkg/
│   └── proto/        # Protocol Buffer definitions
│       ├── greet.proto
│       └── generated/  # Generated Go code from .proto files (will be created)
├── scripts/
│   └── generate-proto.sh # Script to automate protobuf code generation
├── go.mod            # Go module definition
├── go.sum            # Go module checksums
└── README.md

## Getting Started
These instructions will get a copy of the project up and running on your local machine, primarily focusing on a WSL (Windows Subsystem for Linux) Ubuntu environment.

### Prerequisites
Ensure you have the following installed in your WSL Ubuntu terminal:

- Go (1.20 or newer recommended): Verify with go version. If not installed, follow instructions on go.dev/doc/install. A common method is:
Bash
```
wget https://go.dev/dl/go1.22.3.linux-amd64.tar.gz # Check go.dev for latest version
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc # Or close/reopen terminal
```
- Protocol Buffer Compiler (protoc): Verify with which protoc. If not installed, download the protoc-*-linux-x86_64.zip from protobuf releases, then:
Bash
```
mkdir -p /tmp/protoc_install
unzip protoc-26.1-linux-x86_64.zip -d /tmp/protoc_install # Adjust version as needed
sudo mv /tmp/protoc_install/bin/protoc /usr/local/bin/
sudo cp -R /tmp/protoc_install/include/* /usr/local/include/
rm -rf /tmp/protoc_install
```
- Go Protobuf Plugins: These are installed via go install. Crucially, ensure your $HOME/go/bin directory is in your PATH.
```
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.bashrc # Add this if not already there
source ~/.bashrc # Or close/reopen terminal
```
Verify with which protoc-gen-go and which protoc-gen-go-grpc.

grpcurl (for testing):
Bash
```
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```
### Verify with which grpcurl. Ensure $HOME/go/bin is in your PATH (same as Go plugins).

### Installation & Setup
1 - Clone the repository:

Bash
```
git clone [your-repo-url]
cd grpc-greeter
```
(If you're building this from scratch, simply navigate to your project directory grpc-greeter.)

2 - Generate Go Protobuf Code:
This step uses protoc and the Go plugins to generate .pb.go files from your .proto definitions.

Bash
```
cd scripts
./generate-proto.sh
cd .. # Go back to project root
```
You should now see greet.pb.go and greet_grpc.pb.go inside pkg/proto/generated/.

3 - Download Go Modules:

Bash
```
go mod tidy
```
This will download all necessary Go dependencies.

4 - Build the Server Application:

Bash
```
go build -o ./bin/server ./cmd/server
```
This compiles the server executable into a bin directory (which will be created).

5 - Create config.yaml:
The server expects a config.yaml file in the project root (grpc-greeter/). Create this file:

YAML
```
# config.yaml
grpc_port: 50051
```
## Running the Service
Start the gRPC Server:
Open a new WSL terminal window, navigate to your project root (grpc-greeter/), and run the server:

Bash
```
./bin/server
```
You should see output similar to:

INFO: YYYY/MM/DD HH:MM:SS logger.go:42: Starting gRPC server...
INFO: YYYY/MM/DD HH:MM:SS logger.go:42: gRPC server listening on port 50051
Keep this terminal running.

Test with grpcurl (in a separate terminal):
Open another new WSL terminal window, navigate to your project root, and try these commands:

List exposed services:

Bash
```
grpcurl -plaintext localhost:50051 list
```
You should see:

greet.Greeter
grpc.reflection.v1.ServerReflection
grpc.reflection.v1alpha.ServerReflection
Invoke the SayHello method:

Bash
```
grpcurl -plaintext -d '{"name": "World"}' localhost:50051 greet.Greeter/SayHello
```
You should receive a successful response, like:
```
JSON

{
  "message": "Hello, World !"
}
```
(Note: The exact message might vary slightly based on your FormatGreeting implementation).

Run the Go Client (Optional, if you implement it):
If you've implemented cmd/client/main.go, you can build and run it similarly:

Bash
```
go build -o ./bin/client ./cmd/client
./bin/client
```
