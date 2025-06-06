# Factura Móvil Gateway

A Go-based HTTP API service for handling Chilean electronic invoicing (DTE - Documento Tributario Electrónico) operations. This gateway manages electronic invoice authorization files (CAF - Código de Autorización de Folios) and provides document stamping services.

## Features

- **CAF Management**: Handle storage, retrieval, and lifecycle management of electronic invoice authorization files
- **Document Stamping**: Provide digital stamping services for electronic documents
- **Company Management**: Manage company information and their associated authorization files
- **SII Compliance**: Fully compliant with Chilean SII (Servicio de Impuestos Internos) electronic invoicing standards

## Prerequisites

- Go 1.23.6 or later
- Docker and Docker Compose
- PostgreSQL (via Docker)

## Local Development Setup

### 1. Clone the Repository
```bash
git clone <repository-url>
cd fm-gateway
```

### 2. Start PostgreSQL Database
```bash
# Start PostgreSQL using Docker Compose
docker-compose up -d postgres

# Check database status
docker-compose ps

# View database logs
docker logs fm-gateway-postgres
```

### 3. Configure Environment
```bash
# Use local development environment
source .envrc.local

# Or manually set environment variables
export FMG_DBHOST=localhost
export FMG_DBUSER=fmgateway
export FMG_DBPASS=fmgateway123
```

### 4. Install Dependencies
```bash
go mod download
go mod tidy
```

### 5. Run the Application
```bash
# Run the API server
go run cmd/api/main.go
```

The application will start on the default port with the following services:
- HTTP API endpoints
- Prometheus metrics endpoint
- Database connectivity
- File storage in `tmp/` directory

## Database Management

### Basic Operations
```bash
# Start database
docker-compose up -d postgres

# Stop database
docker-compose down

# Restart database
docker-compose restart postgres

# View database logs
docker logs fm-gateway-postgres --follow
```

### Direct Database Access
```bash
# Connect to PostgreSQL
docker exec -it fm-gateway-postgres psql -U fmgateway -d postgres

# Run SQL commands
docker exec -it fm-gateway-postgres psql -U fmgateway -d postgres -c "SELECT version();"
```

### Database Backup and Restore
```bash
# Backup database
docker exec fm-gateway-postgres pg_dump -U fmgateway postgres > backup.sql

# Restore database
docker exec -i fm-gateway-postgres psql -U fmgateway postgres < backup.sql
```

## Development Workflow

### Environment Switching
```bash
# For local development
source .envrc.local

# For remote/production environment
source .envrc
```

### Testing the API
```bash
# Health check (if implemented)
curl http://localhost:8080/health

# Prometheus metrics
curl http://localhost:8080/metrics
```

### Code Quality
```bash
# Format code
gofmt -w .

# Run pre-commit hooks
pre-commit run --all-files

# Run tests (when available)
go test ./...

# Run with race detection
go test -race ./...
```

## Project Structure

```
fm-gateway/
├── cmd/api/              # Application entry point
├── internal/
│   ├── controllers/      # HTTP request handlers
│   ├── datatypes/        # Shared data structures
│   ├── domain/          # Business entities and rules
│   ├── httpserver/      # HTTP server configuration
│   ├── persistence/     # Database repositories
│   ├── storage/         # File storage operations
│   └── usecases/        # Application business logic
├── schemas/             # SII XML schemas
├── tmp/                 # File storage directory
├── docker-compose.yml   # Docker services configuration
├── .envrc              # Production environment variables
├── .envrc.local        # Local development environment variables
└── go.mod              # Go module dependencies
```

## API Endpoints

The application provides REST API endpoints for:
- CAF (Código de Autorización de Folios) operations
- Document stamping services
- Company management

*Note: Detailed API documentation will be added as endpoints are finalized.*

## Configuration

### Environment Variables

| Variable | Description | Local Development | Production |
|----------|-------------|-------------------|------------|
| `FMG_DBHOST` | Database host | `localhost` | Production DB host |
| `FMG_DBUSER` | Database username | `fmgateway` | Production DB user |
| `FMG_DBPASS` | Database password | `fmgateway123` | Production DB password |

### Database Configuration
- **Database**: PostgreSQL 15
- **Port**: 5432
- **Default Database**: `postgres`
- **Connection Pooling**: Managed by GORM

## Monitoring and Observability

### Prometheus Metrics
The application exposes Prometheus metrics for monitoring:
```bash
curl http://localhost:8080/metrics
```

### Logging
- Structured logging using Go's `slog` package
- Debug level logging with source attribution
- Logs written to stdout

### Health Checks
```bash
# PostgreSQL health check
docker exec fm-gateway-postgres pg_isready -U fmgateway -d postgres
```

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   ```bash
   # Check if PostgreSQL is running
   docker-compose ps
   
   # Check database logs
   docker logs fm-gateway-postgres
   
   # Restart database
   docker-compose restart postgres
   ```

2. **Port Already in Use**
   ```bash
   # Check what's using port 5432
   lsof -i :5432
   
   # Stop conflicting services or change port in docker-compose.yml
   ```

3. **Environment Variables Not Set**
   ```bash
   # Verify environment variables
   env | grep FMG_
   
   # Source the correct environment file
   source .envrc.local
   ```

## Contributing

1. Follow Go best practices and the existing code structure
2. Run `gofmt` before committing
3. Ensure all tests pass
4. Update documentation as needed

## License

[Add your license information here] 