# System Patterns - Factura Móvil Gateway

## Architecture Overview

### Clean Architecture Implementation
The system follows Clean Architecture principles with clear separation of concerns:

```
cmd/api/          - Application entry point
internal/
├── domain/       - Business entities and rules
├── usecases/     - Application business logic
├── controllers/  - HTTP request handlers
├── persistence/  - Database repositories
├── storage/      - File storage operations
├── httpserver/   - HTTP server configuration
└── datatypes/    - Shared data structures
```

### Dependency Direction
- **Outer layers depend on inner layers, never the reverse**
- Controllers → Use Cases → Domain
- Persistence implements interfaces defined in use cases
- Domain layer has no external dependencies

## Key Design Patterns

### Repository Pattern
- **Implementation**: `persistence/` package contains repository implementations
- **Interfaces**: Defined in use case layer for dependency inversion
- **Benefits**: Database-agnostic business logic, easy testing with mocks

### Builder Pattern
- **Used in**: CAF domain entity creation (`NewCAFBuilder()`)
- **Benefits**: Complex object construction with validation
- **Pattern**: Fluent interface with method chaining

### Dependency Injection
- **Manual DI**: Dependencies injected in `main.go`
- **No Framework**: Uses constructor functions for clean initialization
- **Testing**: Easy to mock dependencies for unit tests

### Service Layer Pattern
- **Use Cases**: Encapsulate business operations
- **Single Responsibility**: Each service handles one domain area
- **Examples**: `CAFService`, `StampService`, `CompanyService`

## Component Relationships

### Data Flow
1. **HTTP Request** → Controller
2. **Controller** → Use Case Service
3. **Service** → Repository (for data) + Storage (for files)
4. **Repository** → Database
5. **Storage** → File System

### Entity Lifecycle
- **CAF Files**: Upload → Validation → Storage → Database Record → Usage Tracking
- **Companies**: Registration → CAF Association → Monitoring
- **Documents**: Submission → Stamping → Return

## Technical Decisions

### Database Design
- **GORM**: ORM for database operations
- **PostgreSQL**: Chosen for reliability and ACID compliance
- **Migrations**: Handle schema evolution (assumed, not visible in current code)

### File Storage
- **Local Storage**: `tmp/` directory for file operations
- **Abstraction**: Storage interface allows for future cloud storage integration
- **Security**: Files isolated from direct web access

### Error Handling
- **Go Idioms**: Explicit error returns from all operations
- **Error Wrapping**: Using `fmt.Errorf` with `%w` for context
- **No Panics**: Graceful error handling throughout the system

### Configuration
- **Environment Variables**: All configuration through env vars
- **Validation**: Configuration validated at startup
- **Security**: Sensitive data (DB credentials) not hardcoded

### Logging
- **Structured Logging**: Using `slog` package
- **Debug Level**: Comprehensive logging for troubleshooting
- **Source Attribution**: File and line information included

### HTTP Server
- **Graceful Shutdown**: Signal handling for clean termination
- **CORS Support**: Cross-origin requests enabled
- **Metrics**: Prometheus integration for monitoring

## Extensibility Points
- **Storage Interface**: Easy to add cloud storage backends
- **Repository Interfaces**: Can swap database implementations
- **Controller Pattern**: New endpoints follow established patterns
- **Service Layer**: New business operations easily added 