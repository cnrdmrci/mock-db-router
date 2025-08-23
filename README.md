# Mock Database Router

A high-performance HTTP mock server that serves responses from a PostgreSQL database. This application allows you to create dynamic mock APIs by storing request patterns and their corresponding responses in a database.

## 🚀 Features

- **Dynamic Mock Responses**: Store and serve mock responses based on URL path and HTTP method
- **PostgreSQL Integration**: All mock data is stored in PostgreSQL for persistence and easy management
- **HTTP Method Support**: Different responses for GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD
- **Query Parameter Support**: Full URL path including query parameters for precise matching
- **Custom Headers**: Set custom response headers stored as key=value pairs
- **Custom Status Codes**: Return any HTTP status code (200, 404, 500, etc.)
- **High Performance**: Connection pooling for optimal database performance
- **Concurrent Safe**: Handles multiple simultaneous requests efficiently

## 📋 Prerequisites

- Go 1.22.1 or higher
- PostgreSQL database
- Network access to your PostgreSQL instance

## 🛠️ Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/cnrdmrci/mock-db-router
   cd mock-db-router
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up PostgreSQL database**
   
   Create the required table:
   ```sql
   CREATE TABLE IF NOT EXISTS public.mock_responses (
       id SERIAL PRIMARY KEY,
       path VARCHAR(500) NOT NULL,
       method VARCHAR(10) NOT NULL,
       request_body JSONB,
       response_body JSONB NOT NULL,
       response_status_code INTEGER DEFAULT 200,
       headers TEXT,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );
   ```

4. **Configure database connection**
   
   Update the connection string in `main.go`:
   ```go
   const connStr = "host=localhost port=5432 user=your_user password=your_password dbname=your_db sslmode=disable"
   ```

## 🏃‍♂️ Usage

### Starting the Server

```bash
go run main.go
```

The server will start on port 8080 by default.

### Adding Mock Responses

Insert mock responses into the database:

```sql
-- Example: GET request
INSERT INTO mock_responses (path, method, response_body, headers, response_status_code) 
VALUES (
    '/api/users/123', 
    'GET', 
    '{"id": 123, "name": "John Doe", "email": "john@example.com"}', 
    'Content-Type=application/json;Cache-Control=no-cache', 
    200
);

-- Example: POST request
INSERT INTO mock_responses (path, method, request_body, response_body, headers, response_status_code) 
VALUES (
    '/api/users', 
    'POST',
    '{"name": "John Doe", "email": "john@example.com"}', 
    '{"message": "User created successfully", "id": 124}', 
    'Content-Type=application/json;Location=/api/users/124', 
    201
);

-- Example: With query parameters
INSERT INTO mock_responses (path, method, response_body, headers, response_status_code) 
VALUES (
    '/api/users?active=true&page=1', 
    'GET', 
    '[{"id": 1, "name": "Active User 1"}, {"id": 2, "name": "Active User 2"}]', 
    'Content-Type=application/json', 
    200
);

-- Example: Error response
INSERT INTO mock_responses (path, method, response_body, headers, response_status_code) 
VALUES (
    '/api/users/999', 
    'GET', 
    '{"error": "User not found", "code": "USER_NOT_FOUND"}', 
    'Content-Type=application/json', 
    404
);
```

### Making Requests

Once the server is running and you have mock data in the database:

```bash
# GET request
curl -X GET http://localhost:8080/api/users/123

# POST request
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name": "New User"}'

# Request with query parameters
curl -X GET "http://localhost:8080/api/users?active=true&page=1"
```

## 📊 Database Schema

### Table: `mock_responses`

| Column | Type | Description |
|--------|------|-------------|
| `id` | SERIAL | Primary key |
| `path` | VARCHAR(500) | Full URL path including query parameters |
| `method` | VARCHAR(10) | HTTP method (GET, POST, PUT, DELETE, etc.) |
| `request_body` | JSONB | Request body content |
| `response_body` | JSONB | Response content to return |
| `response_status_code` | INTEGER | HTTP status code (default: 200) |
| `headers` | TEXT | Headers in "key=value;key2=value2" format |
| `created_at` | TIMESTAMP | Record creation timestamp |

## ⚙️ Configuration

### Database Connection Pool

The application uses connection pooling for optimal performance:

- **Max Open Connections**: 10
- **Max Idle Connections**: 5
- **Connection Lifetime**: 15 minutes
- **Idle Timeout**: 3 minutes

### Headers Format

Headers should be stored as semicolon-separated key=value pairs:
