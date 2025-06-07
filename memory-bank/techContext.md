# Technical Context - Factura Móvil Gateway

## Core Technologies

### Language & Runtime
- **Go 1.23.6**: Primary development language
- **Module**: `factura-movil-gateway` with semantic versioning

### Database & Persistence
- **PostgreSQL**: Primary database for structured data
- **GORM 1.26.1**: ORM for database operations
  - `gorm.io/driver/postgres v1.5.11`: PostgreSQL driver
  - Migration support and relationship management

### HTTP & API
- **Standard Library**: `net/http` for HTTP server operations
- **CORS**: `github.com/rs/cors v1.11.1` for cross-origin requests
- **REST Architecture**: JSON-based API endpoints

### File & Storage Management
- **Local File System**: File storage in `tmp/` directory
- **Storage Abstraction**: Interface-based design for future extensibility

### Monitoring & Observability
- **Prometheus**: `github.com/prometheus/client_golang v1.22.0`
  - Metrics collection and monitoring
  - Standard Go application metrics
- **Structured Logging**: Go's `log/slog` package
  - Debug-level logging with source attribution
  - JSON and text output formats

### Utilities
- **UUID Generation**: `github.com/google/uuid v1.6.0`
- **XML Processing**: Standard library `encoding/xml`
- **Context Management**: Standard library `context` for cancellation

## Development Environment

### Prerequisites
- Go 1.23.6 or later
- PostgreSQL database server
- Environment variable management (`.envrc` suggests direnv usage)

### Environment Configuration
Required environment variables:
- `FMG_DBHOST`: Database host address
- `FMG_DBUSER`: Database username
- `FMG_DBPASS`: Database password

### Project Structure Standards
- **Clean Architecture**: Enforced directory structure
- **Internal Package**: All business logic in `internal/`
- **Command Pattern**: Entry points in `cmd/`
- **Schema Validation**: XSD files in `schemas/`

### Development Tools
- **Just**: Task runner for build and run commands (`justfile`)
- **Pre-commit Hooks**: `.pre-commit-config.yaml` for code quality
- **Git Integration**: Standard Git workflow
- **Environment Management**: `direnv` with `.envrc`

## Dependencies Deep Dive

### Direct Dependencies
```go
require (
    github.com/google/uuid v1.6.0
    github.com/prometheus/client_golang v1.22.0
    github.com/rs/cors v1.11.1
    gorm.io/driver/postgres v1.5.11
    gorm.io/gorm v1.26.1
)
```

### Transitive Dependencies
- **Database**: PostgreSQL connection pooling, prepared statements
- **Cryptography**: TLS support, password hashing
- **Network**: HTTP/2 support, connection management
- **Monitoring**: Prometheus exposition format

## Compliance & Standards

### Chilean SII Standards
- **DTE v1.0**: Electronic document schemas
- **XML Digital Signature**: Document authentication
- **CAF Format**: Authorization file structure
- **SII Types**: Chilean tax authority data types

### Schema Files
- `DTE_v10.xsd`: Primary electronic document schema (227KB, 5321 lines)
- `EnvioDTE_v10.xsd`: Document transmission schema
- `SiiTypes_v10.xsd`: Chilean tax system data types
- `xmldsignature_v10.xsd`: Digital signature standards

## Performance Considerations
- **Connection Pooling**: GORM handles database connection management
- **Concurrent Safety**: Go routines with proper synchronization
- **Memory Management**: Efficient XML processing for large documents
- **File I/O**: Streaming for large CAF files

## Security Architecture
- **Environment-based Configuration**: No hardcoded secrets
- **Database Security**: Parameterized queries via GORM
- **File System Security**: Isolated tmp directory
- **TLS Support**: HTTPS endpoint capability
- **Input Validation**: XML schema validation

## Development Workflow

### Build & Run Commands
- **Run Application**: `just run` (equivalent to `go run ./cmd/api`)
- **Build Application**: `just build` (creates binary at `bin/fm-gateway`)
- **Run Tests**: `just test` (standard Go tests)
- **Development Mode**: `just dev` (live reload with air)
- **Code Formatting**: `just fmt` (go fmt)
- **Linting**: `just lint` (golangci-lint)
- **Dependencies**: `just deps` (go mod operations)

**Note**: Always use `just` commands instead of direct `go` commands for consistency and workflow standardization.

## Deployment Architecture
- **Single Binary**: Compiled Go application
- **External Dependencies**: PostgreSQL database
- **Storage Requirements**: Local filesystem access
- **Monitoring Integration**: Prometheus metrics endpoint
- **Graceful Shutdown**: Signal handling for clean termination 