# API Documentation - Factura Móvil Gateway

## Overview

The Factura Móvil Gateway provides a complete RESTful API for managing Chilean electronic invoicing operations. The API is designed around a company-centric architecture where Companies are the primary resource, and CAFs (Código de Autorización de Folios) and Stamps are sub-resources.

## Base URL

```
http://localhost:8080
```

## Authentication

Currently, the API does not implement authentication. This should be added before production deployment.

## API Endpoints

### Company Management

#### Create Company
Create a new company in the system.

**Endpoint:** `POST /companies`

**Request Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Empresa Ejemplo S.A.",
  "code": "12345678-9",
  "factura_movil_company_id": 12345
}
```

**Response:**
- **Status:** `201 Created`
- **Body:**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "Empresa Ejemplo S.A.",
  "code": "12345678-9",
  "factura_movil_company_id": 12345
}
```

**Error Responses:**
- `400 Bad Request`: Invalid JSON or missing required fields
- `500 Internal Server Error`: Database or server error

---

#### List Companies
Retrieve a list of all companies with optional filtering.

**Endpoint:** `GET /companies`

**Query Parameters:**
- `name` (optional): Filter companies by partial name match (case-insensitive)

**Examples:**
```
GET /companies
GET /companies?name=ejemplo
GET /companies?name=corp
```

**Response:**
- **Status:** `200 OK`
- **Body:**
```json
[
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "name": "Empresa Ejemplo S.A.",
    "code": "12345678-9",
    "factura_movil_company_id": 12345
  },
  {
    "id": "456e7890-e89b-12d3-a456-426614174000",
    "name": "Otra Empresa Corp.",
    "code": "87654321-0",
    "factura_movil_company_id": 67890
  }
]
```

**Error Responses:**
- `500 Internal Server Error`: Database or server error

---

#### Get Company by ID
Retrieve a specific company by its ID.

**Endpoint:** `GET /companies/{companyId}`

**Path Parameters:**
- `companyId`: UUID of the company

**Example:**
```
GET /companies/123e4567-e89b-12d3-a456-426614174000
```

**Response:**
- **Status:** `200 OK`
- **Body:**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "Empresa Ejemplo S.A.",
  "code": "12345678-9",
  "factura_movil_company_id": 12345
}
```

**Error Responses:**
- `404 Not Found`: Company with specified ID not found
- `500 Internal Server Error`: Database or server error

---

### CAF Management

#### Upload CAF for Company
Upload and process a CAF (Código de Autorización de Folios) file for a specific company.

**Endpoint:** `POST /companies/{companyId}/cafs`

**Path Parameters:**
- `companyId`: UUID of the company

**Request Headers:**
```
Content-Type: application/xml
```

**Request Body:**
CAF XML file content (ISO-8859-1 encoding supported)

```xml
<AUTORIZACION>
  <CAF>
    <DA>
      <RE>12345678-9</RE>
      <RS>Empresa Ejemplo S.A.</RS>
      <TD>33</TD>
      <RNG>
        <D>1</D>
        <H>100</H>
      </RNG>
      <FA>2024-01-15</FA>
      <!-- Additional CAF XML content -->
    </DA>
  </CAF>
</AUTORIZACION>
```

**Response:**
- **Status:** `201 Created`
- **Body:**
```json
{
  "id": "caf-uuid-generated",
  "companyId": "123e4567-e89b-12d3-a456-426614174000",
  "companyCode": "12345678-9",
  "companyName": "Empresa Ejemplo S.A.",
  "documentType": 33,
  "initialFolios": 1,
  "currentFolios": 1,
  "finalFolios": 100,
  "authorizationDate": "2024-01-15T00:00:00Z",
  "expirationDate": "2024-07-15T00:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid XML format or CAF data
- `404 Not Found`: Company not found
- `500 Internal Server Error`: Database, storage, or server error

**File Storage:**
The CAF XML file is stored at: `tmp/caf/{companyId}/{cafId}.xml`

---

#### List CAFs for Company
Retrieve all CAFs associated with a specific company.

**Endpoint:** `GET /companies/{companyId}/cafs`

**Path Parameters:**
- `companyId`: UUID of the company

**Example:**
```
GET /companies/123e4567-e89b-12d3-a456-426614174000/cafs
```

**Response:**
- **Status:** `200 OK`
- **Body:**
```json
[
  {
    "id": "caf-uuid-1",
    "companyId": "123e4567-e89b-12d3-a456-426614174000",
    "companyCode": "12345678-9",
    "companyName": "Empresa Ejemplo S.A.",
    "documentType": 33,
    "initialFolios": 1,
    "currentFolios": 1,
    "finalFolios": 100,
    "authorizationDate": "2024-01-15T00:00:00Z",
    "expirationDate": "2024-07-15T00:00:00Z"
  },
  {
    "id": "caf-uuid-2",
    "companyId": "123e4567-e89b-12d3-a456-426614174000",
    "companyCode": "12345678-9",
    "companyName": "Empresa Ejemplo S.A.",
    "documentType": 34,
    "initialFolios": 101,
    "currentFolios": 101,
    "finalFolios": 200,
    "authorizationDate": "2024-01-20T00:00:00Z",
    "expirationDate": "2024-07-20T00:00:00Z"
  }
]
```

**Error Responses:**
- `404 Not Found`: Company not found
- `500 Internal Server Error`: Database or server error

---

### Document Stamping

#### Generate Stamp for Company
Generate a digital stamp (timbre) for a document using company-specific data.

**Endpoint:** `POST /companies/{companyId}/stamps`

**Path Parameters:**
- `companyId`: UUID of the company

**Request Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "fmaPago": "1",
  "hasTaxes": true,
  "details": [
    {
      "position": 1,
      "product": {
        "unit": {
          "code": "UN"
        },
        "price": 10000,
        "name": "Producto Ejemplo",
        "code": "PROD001"
      },
      "description": "Descripción del producto",
      "quantity": 2.0,
      "discount": 0.0
    }
  ],
  "client": {
    "address": "Dirección Cliente 123",
    "name": "Cliente Ejemplo",
    "municipality": "Santiago",
    "line": "Giro Cliente",
    "code": "11111111-1"
  },
  "assignedFolio": "1",
  "subsidiary": {
    "code": "001"
  },
  "date": "2024-01-15"
}
```

**Response:**
- **Status:** `200 OK`
- **Content-Type:** `application/xml`
- **Body:**
```xml
<TED version="1.0">
  <DD>
    <RE>12345678-9</RE>
    <TD>33</TD>
    <F>1</F>
    <FE>2024-01-15</FE>
    <RR>11111111-1</RR>
    <RSR>Cliente Ejemplo</RSR>
    <MNT>20000</MNT>
    <IT1>Producto Ejemplo</IT1>
    <CAF>FIXME</CAF>
    <TSTED>2024-01-15T10:30:00-03:00</TSTED>
  </DD>
  <FRMT></FRMT>
</TED>
```

**Field Descriptions:**
- `RE`: Company code (automatically filled from company data)
- `TD`: Document type (33 for invoices with taxes, 34 without taxes)
- `F`: Folio number
- `FE`: Document creation date
- `RR`: Customer code
- `RSR`: Customer name
- `MNT`: Total amount
- `IT1`: First item name
- `CAF`: CAF reference (placeholder)
- `TSTED`: Timestamp when stamp was generated

**Error Responses:**
- `400 Bad Request`: Invalid JSON or invoice data
- `404 Not Found`: Company not found
- `500 Internal Server Error`: Stamp generation or server error

---

## Error Response Format

All error responses follow a consistent format:

**Error Response:**
```json
{
  "error": "Description of the error"
}
```

**Common HTTP Status Codes:**
- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request data
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server-side error

## Data Models

### Company
```json
{
  "id": "string (UUID)",
  "name": "string",
  "code": "string (RUT format: 12345678-9)",
  "factura_movil_company_id": "integer"
}
```

### CAF
```json
{
  "id": "string (UUID)",
  "companyId": "string (UUID, foreign key)",
  "companyCode": "string (RUT from XML)",
  "companyName": "string",
  "documentType": "integer (33=with taxes, 34=without taxes)",
  "initialFolios": "integer",
  "currentFolios": "integer",
  "finalFolios": "integer",
  "authorizationDate": "string (ISO 8601 timestamp)",
  "expirationDate": "string (ISO 8601 timestamp, +6 months from authorization)"
}
```

## Usage Examples

### Complete Workflow Example

1. **Create a Company:**
```bash
curl -X POST http://localhost:8080/companies \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Mi Empresa S.A.",
    "code": "76543210-9",
    "factura_movil_company_id": 54321
  }'
```

2. **Upload a CAF for the Company:**
```bash
curl -X POST http://localhost:8080/companies/company-uuid/cafs \
  -H "Content-Type: application/xml" \
  -d @caf-file.xml
```

3. **Generate a Stamp:**
```bash
curl -X POST http://localhost:8080/companies/company-uuid/stamps \
  -H "Content-Type: application/json" \
  -d '{
    "hasTaxes": true,
    "client": {"code": "11111111-1", "name": "Cliente"},
    "details": [{"position": 1, "product": {"name": "Producto", "price": 1000}, "quantity": 1}],
    "date": "2024-01-15"
  }'
```

4. **List Company CAFs:**
```bash
curl http://localhost:8080/companies/company-uuid/cafs
```

## Rate Limiting

Currently, no rate limiting is implemented. This should be added for production use.

## Versioning

The API currently doesn't implement versioning. Future versions should include version headers or URL versioning.

## Security Considerations

1. **Authentication**: Not implemented - should be added before production
2. **Input Validation**: Basic validation implemented, but should be enhanced
3. **Rate Limiting**: Not implemented - recommended for production
4. **HTTPS**: Should be enforced in production
5. **CORS**: Basic CORS support implemented

## Development and Testing

### Environment Variables
```bash
export FMG_DBHOST=localhost
export FMG_DBUSER=postgres
export FMG_DBPASS=password
```

### Starting the Server
```bash
go run cmd/api/main.go
```

The server will start on port 8080 with logging enabled. 