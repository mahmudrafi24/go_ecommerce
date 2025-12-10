# Go E-commerce Demo

A simple e-commerce backend service written in Go.

## Features

*   User authentication
*   Product management
*   Shopping cart functionality
*   Order processing

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

*   [Go](https://golang.org/doc/install) (version 1.22 or later)
*   [PostgreSQL](https://www.postgresql.org/download/)
*   [Docker](https://www.docker.com/products/docker-desktop) (optional, for running PostgreSQL)

### Installation

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/your-username/go-ecommerce-demo.git
    cd go-ecommerce-demo
    ```

2.  **Set up environment variables:**

    Create a `.env` file in the root of the project and add the following variables:

    ```
    PORT=8080
    DATABASE_URL="postgres://youruser:yourpassword@localhost:5432/ecommerce?sslmode=disable"
    JWT_SECRET="your-super-secret-key"
    ```

    **Note:** You will need to create a PostgreSQL database named `ecommerce` and update the `DATABASE_URL` with your credentials.

3.  **Install dependencies:**

    ```bash
    go mod tidy
    ```

## Usage

To run the application, execute the following command from the root directory:

```bash
go run cmd/server/main.go
```

The server will start on the port specified in your `.env` file (defaulting to `8080`).

## API Endpoints

The following are the available API endpoints.

### Health Check

*   **GET /health**

    Returns the status of the server.

    **Response:**

    ```json
    {
      "status": "ok"
    }
    ```

*(Further endpoints for users, products, cart, and orders will be documented here.)*

## Project Structure

```
.
├── cmd
│   └── server
│       └── main.go         # Application entry point
├── internal
│   ├── config              # Configuration management
│   ├── database            # Database connection and setup
│   ├── handler             # HTTP handlers
│   ├── middleware          # Custom middleware
│   ├── models              # Data models and structures
│   ├── repository          # Database repository layer
│   └── service             # Business logic
├── migrations              # Database migrations
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
└── README.md               # This file
```

## Dependencies

*   [chi](github.com/go-chi/chi/v5) - Lightweight, idiomatic and composable router for building Go HTTP services.
*   [pgx](github.com/jackc/pgx/v5) - PostgreSQL driver and toolkit for Go.
