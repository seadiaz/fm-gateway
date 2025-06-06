# System Patterns - Factura Móvil Gateway

## Architecture Overview

### Clean Architecture Implementation
The system successfully implements Clean Architecture principles with complete separation of concerns and proper dependency direction:

```
cmd/api/          - Application entry point with dependency injection
internal/
├── domain/       - Business entities (Company, CAF) with business rules
├── usecases/     - Application services with business logic
├── controllers/  - HTTP request handlers with validation
├── persistence/  - Database repositories with GORM
├── storage/      - File storage operations (company-organized)
├── httpserver/   - HTTP server configuration and utilities
└── datatypes/    - Shared data structures and XML types
```

### Dependency Direction Enforced
- **✅ Controllers → Use Cases → Domain**: All dependencies point inward
- **✅ Persistence implements interfaces**: Defined in use case layer for dependency inversion
- **✅ Domain layer isolation**: No external dependencies in business logic
- **✅ Clean boundaries**: Each layer has clearly defined responsibilities

## Key Design Patterns Implemented

### Repository Pattern
- **✅ Complete Implementation**: All repositories follow consistent patterns
  - `CompanyRepository`: FindAll, FindByID, FindByNameFilter, Save
  - `CAFRepository`: Save, FindByCompanyID with company association
- **✅ Interface Abstractions**: Defined in use case layer for testability
- **✅ GORM Integration**: Database-agnostic business logic with ORM
- **✅ Error Handling**: Consistent error wrapping and context support

### Builder Pattern
- **✅ Domain Entity Construction**: Complex object creation with validation
  - `CompanyBuilder`: WithName, WithCode, WithFacturaMovilCompanyID
  - `CAFBuilder`: Enhanced with WithCompanyID and WithCompanyCode
- **✅ Fluent Interface**: Method chaining for readable construction
- **✅ Validation**: Build-time validation ensures entity consistency

### Service Layer Pattern
- **✅ Business Logic Encapsulation**: Clear separation of concerns
  - `CompanyService`: Company management with filtering capabilities
  - `CAFService`: CAF lifecycle with company association
  - `StampService`: Dynamic stamp generation using company data
- **✅ Single Responsibility**: Each service handles one domain area
- **✅ Interface-Based**: Easy to mock and test

### Dependency Injection Pattern
- **✅ Constructor Injection**: Manual DI in main.go for simplicity
- **✅ Interface Dependencies**: Services depend on abstractions
- **✅ No Framework Dependency**: Clean, testable dependency management
- **✅ Explicit Wiring**: Clear dependency relationships in main.go

## RESTful API Patterns

### Resource-Oriented Design
- **✅ Primary Resources**: Companies as main business entities
- **✅ Sub-Resources**: CAFs and Stamps as company sub-resources
- **✅ HTTP Methods**: Proper use of GET, POST for different operations
- **✅ Status Codes**: Appropriate HTTP responses (200, 201, 400, 404, 500)

### Resource Hierarchy
```
/companies                     # Company collection
├── POST   /companies          # Create company
├── GET    /companies          # List companies (with filtering)
├── GET    /companies/{id}     # Get specific company
├── POST   /companies/{id}/cafs     # Upload CAF for company
├── GET    /companies/{id}/cafs     # List company CAFs
└── POST   /companies/{id}/stamps   # Generate company stamps
```

### Controller Pattern
- **✅ Single Responsibility**: Each controller handles one resource type
- **✅ HTTP Handler Functions**: Clean handler function organization
- **✅ Path Parameter Extraction**: Company ID validation in sub-resources
- **✅ Consistent Error Handling**: Standardized error responses

## Data Layer Patterns

### Entity Relationship Design
```sql
-- Primary Entity
Company (id, name, code, factura_movil_company_id)
    ↓ (1:N relationship)
-- Associated Entity with dual identification
CAF (id, company_id, company_code, company_name, ...)
```

### Dual Identification Pattern
- **CompanyID**: Database foreign key relationship (UUID)
- **CompanyCode**: Business identifier from XML (e.g., "12345678-9")
- **Separation of Concerns**: Technical vs business identification

### Database Optimization Patterns
- **✅ Indexing Strategy**: Index on company_id for optimal CAF queries
- **✅ Repository Queries**: Optimized queries with proper WHERE clauses
- **✅ Context Support**: All database operations use context.Context
- **✅ Error Mapping**: GORM errors mapped to domain errors

## Storage Organization Patterns

### Hierarchical Storage Pattern
```
tmp/
├── caf/
│   ├── {companyId}/          # Company-based organization
│   │   ├── {cafId1}.xml      # Individual CAF files
│   │   └── {cafId2}.xml
│   └── {anotherCompanyId}/
│       └── {cafId3}.xml
```

### Storage Abstraction Pattern
- **✅ Interface-Based**: `BlobStorageClient` interface for multiple implementations
- **✅ Local Implementation**: File system storage with company organization
- **✅ Cloud-Ready**: Interface allows for S3/Azure/GCP implementations
- **✅ Context Support**: All storage operations support cancellation

## Error Handling Patterns

### Consistent Error Response Pattern
```go
// Controller Level
httpserver.ReplyWithError(w, http.StatusNotFound, "company not found")

// Service Level
return fmt.Errorf("finding company by id: %w", err)

// Repository Level
return fmt.Errorf("saving company: %w", err)
```

### Error Propagation Strategy
- **✅ Error Wrapping**: `fmt.Errorf` with `%w` for context preservation
- **✅ HTTP Status Mapping**: Business errors mapped to appropriate HTTP codes
- **✅ Logging Strategy**: Structured logging with context at each layer
- **✅ Client-Friendly Messages**: Generic error messages to clients

## Business Logic Patterns

### Domain-Driven Design Elements
- **✅ Domain Entities**: Company and CAF with business rules
- **✅ Value Objects**: Folio ranges, authorization dates with business logic
- **✅ Business Rules**: 6-month CAF expiration, UUID generation
- **✅ Entity Relationships**: Proper association between Company and CAF

### Validation Patterns
- **✅ Input Validation**: XML parsing with schema validation
- **✅ Business Validation**: Company existence before CAF/Stamp operations
- **✅ Data Integrity**: Foreign key constraints in database
- **✅ Error Boundaries**: Validation at controller and domain levels

## Configuration Patterns

### Environment-Based Configuration
```go
// Environment variable configuration
dbhost := os.Getenv("FMG_DBHOST")
dbuser := os.Getenv("FMG_DBUSER")
dbpass := os.Getenv("FMG_DBPASS")

// Connection string construction
dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=5432 sslmode=disable", 
                   dbhost, dbuser, dbpass)
```

### Dependency Construction Pattern
```go
// Layered dependency construction in main.go
storage := storage.NewLocalStorage("tmp")
repository := persistence.NewCompanyRepository(dsn)
service := usecases.NewCompanyService(repository)
controller := controllers.NewCompanyController(service)
```

## Concurrency Patterns

### Context Propagation
- **✅ Context.Context**: All operations accept context for cancellation
- **✅ HTTP Context**: Request context propagated through all layers
- **✅ Database Context**: GORM operations use request context
- **✅ Storage Context**: File operations support context cancellation

### Thread Safety Considerations
- **Repository Layer**: GORM handles connection pooling and thread safety
- **Service Layer**: Stateless services ensure thread safety
- **Controller Layer**: Request-scoped handlers prevent race conditions

## Monitoring and Observability Patterns

### Structured Logging Pattern
```go
slog.Error("failed to find company", 
           slog.String("Error", err.Error()), 
           slog.String("companyId", companyId))
```

### Metrics Integration
- **✅ Prometheus Ready**: HTTP server configured for metrics collection
- **✅ Graceful Shutdown**: Proper signal handling for clean termination
- **✅ Health Check Ready**: Infrastructure prepared for health endpoints

## Testing Patterns (Ready for Implementation)

### Test Structure Preparation
- **Repository Tests**: Database integration testing with test containers
- **Service Tests**: Business logic testing with mocked repositories
- **Controller Tests**: HTTP testing with mocked services
- **Integration Tests**: End-to-end API testing

### Dependency Injection for Testing
- **Interface-Based**: All dependencies use interfaces for easy mocking
- **Constructor Injection**: Simple dependency replacement in tests
- **No Global State**: All state passed through dependencies

## Security Patterns (Implemented Foundation)

### Input Validation
- **✅ XML Parsing**: Safe XML parsing with charset detection
- **✅ Path Parameter Validation**: Company ID validation in all sub-resource operations
- **✅ JSON Validation**: Request body validation with proper error handling
- **✅ Type Safety**: Strong typing throughout the application

### Error Information Disclosure
- **✅ Generic Error Messages**: Internal errors not exposed to clients
- **✅ Detailed Logging**: Full error context logged server-side
- **✅ HTTP Status Codes**: Appropriate status codes without implementation details

## Extension Points for Future Development

### Service Interfaces
- **Authentication Service**: Ready for pluggable auth implementation
- **Notification Service**: Event-driven notifications for CAF lifecycle
- **Audit Service**: Activity logging and compliance tracking

### Storage Implementations
- **Cloud Storage**: S3, Azure Blob, GCP Storage implementations
- **Database Storage**: Store CAF content in database instead of files
- **Hybrid Storage**: Metadata in database, files in cloud storage

### API Extensions
- **Bulk Operations**: Batch CAF upload and processing
- **Search and Filtering**: Advanced query capabilities
- **Reporting**: CAF usage statistics and company reports
- **Webhooks**: Event notifications for external systems

## Performance Optimization Patterns

### Database Optimization
- **✅ Indexing**: Strategic indexes on foreign keys
- **✅ Query Optimization**: Efficient queries with proper WHERE clauses
- **Connection Pooling**: GORM-managed connection pooling
- **Prepared Statements**: GORM automatically uses prepared statements

### Caching Opportunities
- **Company Data**: Frequently accessed company information
- **CAF Metadata**: Metadata caching for faster retrieval
- **Static Content**: Schema files and configuration data

This architecture provides a solid foundation for Chilean electronic invoicing operations with clear patterns for extension, testing, and production deployment. 