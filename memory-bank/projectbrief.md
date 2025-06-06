# Project Brief - Factura Móvil Gateway

## Project Overview
The Factura Móvil Gateway (fm-gateway) is a Go-based HTTP API service designed to handle Chilean electronic invoicing (DTE - Documento Tributario Electrónico) operations. This gateway serves as an intermediary service for managing electronic invoice authorization files (CAF - Código de Autorización de Folios) and document stamping operations.

## Core Purpose
- **CAF Management**: Handle the storage, retrieval, and lifecycle management of electronic invoice authorization files
- **Document Stamping**: Provide digital stamping services for electronic documents
- **Company Management**: Manage company information and their associated authorization files
- **Integration Gateway**: Act as a reliable API gateway for electronic invoicing operations in Chile

## Key Requirements
1. **Compliance**: Must comply with Chilean SII (Servicio de Impuestos Internos) electronic invoicing standards
2. **Reliability**: Handle critical business operations with proper error handling and logging
3. **Data Persistence**: Maintain CAF files and company data in PostgreSQL database
4. **File Storage**: Store authorization files and documents in local storage
5. **API Design**: Provide clean HTTP REST API endpoints for client integration
6. **Monitoring**: Include Prometheus metrics for observability

## Technical Constraints
- Must use Go 1.23.6+ for development
- PostgreSQL database for persistent data storage
- Local file system for document storage
- Must handle XML schema validation (SII DTE schemas)
- Environment-based configuration (no hardcoded values)

## Success Criteria
- Successfully process CAF authorization files
- Provide reliable document stamping services
- Maintain data integrity and proper error handling
- Enable seamless integration with client applications
- Meet Chilean electronic invoicing compliance requirements 