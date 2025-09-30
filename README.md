# swe-cli - CloudWeGo Microservice Project Generator

A CLI tool for generating CloudWeGo-based microservice projects using Hertz (HTTP framework) and Kitex (RPC framework).

## Features

- Generates CloudWeGo microservice projects with Hertz and Kitex
- Uses sweets-layout as the project template
- Supports custom module names
- Automatically installs dependencies
- Includes Wire dependency injection setup
- GORM for database operations with Goose migrations
- Protobuf validation support

## Installation

### From Source
```bash
# Clone the cli repository
git clone https://github.com/go-sweets/cli
cd cli

# Build the CLI tool
go build -o swe-cli main.go

# Install globally (optional)
cp swe-cli /usr/local/bin/
```

### Using Go Install (when published)
```go
go install github.com/go-sweets/cli@latest
```

## Usage

### Generate a new microservice project
```bash
# Generate project with default module name
swe-cli new helloservice

# Generate project with custom module name
swe-cli new helloservice github.com/myorg/helloservice
```

### Project structure generated
```
helloservice/
├── api/                # Protocol buffer definitions
├── cmd/                # Application entry points
├── internal/           # Internal application code
│   ├── boundedcontexts/  # DDD bounded contexts
│   ├── config/         # Configuration
│   ├── server/         # HTTP and gRPC servers
│   └── service/        # Service implementations
├── etc/                # Configuration files
├── Makefile           # Build and development commands
└── go.mod             # Go module definition
```

### Next steps after generation
```bash
cd helloservice
make init    # Install tools and dependencies
make api     # Generate protobuf code
make gen     # Generate Wire dependency injection
make run     # Run the service
```

## Command Reference

- `swe-cli new <project-name> [module-name]` - Generate a new CloudWeGo microservice project
- `swe-cli version` - Show CLI version
- `swe-cli help` - Show help information

## Requirements

- Go 1.21 or later
- Access to go-sweets repository (for template)

## Template Features

The generated project includes:

- **HTTP Server**: CloudWeGo Hertz with middleware support
- **RPC Server**: CloudWeGo Kitex with interceptors
- **Database**: GORM with MySQL driver and Goose migrations
- **Caching**: Redis integration
- **Dependency Injection**: Google Wire
- **Configuration**: Viper-based configuration management
- **Validation**: Protobuf validation with protoc-gen-validate
- **Monitoring**: Prometheus metrics and tracing support

## License
Apache License Version 2.0, http://www.apache.org/licenses/

