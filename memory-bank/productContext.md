# Product Context - Factura Móvil Gateway

## Why This Project Exists

### The Chilean Electronic Invoicing Challenge
In Chile, all businesses must comply with SII (Servicio de Impuestos Internos) electronic invoicing regulations. This involves:
- Managing CAF (Código de Autorización de Folios) files that authorize invoice number ranges
- Ensuring proper digital stamping of electronic documents
- Maintaining compliance with complex XML schemas and validation rules
- Handling the lifecycle of authorization files (expiration, renewal, etc.)

### The Problem We Solve
Many businesses and software providers need a reliable, compliant gateway service that can:
1. **Simplify CAF Management**: Remove the complexity of directly handling authorization files
2. **Ensure Compliance**: Guarantee that all operations meet SII requirements
3. **Provide Reliability**: Offer a stable service for critical business operations
4. **Enable Integration**: Allow easy integration with existing business systems

## How It Should Work

### User Experience Goals
- **Developers**: Simple REST API that abstracts away SII complexity
- **Businesses**: Reliable service that ensures compliance without manual intervention
- **System Administrators**: Clear monitoring and logging for operational visibility

### Core Workflows
1. **CAF Upload & Processing**
   - Companies upload their CAF authorization files
   - System validates and stores the files securely
   - Automatic tracking of folio usage and expiration dates

2. **Document Stamping**
   - Electronic documents are submitted for digital stamping
   - System applies proper cryptographic signatures
   - Returns stamped documents ready for SII submission

3. **Company Management**
   - Maintain company profiles and their associated CAF files
   - Track usage patterns and upcoming expirations
   - Enable bulk operations for enterprise clients

### Integration Philosophy
- **API-First**: Every operation available through clean REST endpoints
- **Stateless**: No client session management required
- **Idempotent**: Safe to retry operations without side effects
- **Observable**: Comprehensive logging and metrics for monitoring

## Value Proposition
- **Compliance Assurance**: Never worry about SII regulation changes
- **Reduced Complexity**: Focus on business logic, not invoicing infrastructure
- **Operational Reliability**: Built-in error handling and monitoring
- **Developer Friendly**: Clear API documentation and predictable behavior 