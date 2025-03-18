# Building a GraphQL E-commerce API with Go

This guide walks through creating a GraphQL API for a simple e-commerce system using Go. The implementation uses:
- github.com/graphql-go/graphql for implementing GraphQL
- github.com/graphql-go/handler for HTTP handling
- Go's standard HTTP server
- SQLite for the database

Our e-commerce system has the following entities:
- Users
- Products
- Orders
- Order Items (items within an order)

## Table of Contents
1. [Project Structure](#1-project-structure)
2. [Database Implementation](#2-database-implementation)
3. [GraphQL Implementation](#3-graphql-implementation)
4. [Running the Server](#4-running-the-server)
5. [Testing the API](#5-testing-the-api)

## 1. Project Structure

The project is organized with the following structure:

```
go-graphql-ecom/
├── src/
│   ├── api/
│   │   └── main.go           # Main application entry point
│   ├── database/
│   │   ├── db.go             # Database connection and initialization
│   │   └── models.go         # Data models and database operations
│   ├── graphql/
│   │   ├── handler.go        # HTTP handler for GraphQL requests
│   │   ├── resolvers.go      # GraphQL resolver functions
│   │   ├── schema.go         # GraphQL schema definition
│   │   └── types.go          # GraphQL type definitions
│   └── test/
│       └── graphql.http      # HTTP test requests
└── data/
    └── ecommerce.db          # SQLite database file
```

## 2. Database Implementation

The database package (`database/`) handles:

- Database connection using a singleton pattern
- Table creation for users, products, orders, and order items
- CRUD operations for all entities

Key features:
- Thread-safe database initialization with `sync.Once`
- Proper connection management
- Comprehensive data models with relationships

## 3. GraphQL Implementation

The GraphQL implementation is split into several components:

### Types (`types.go`)
Defines GraphQL object types that correspond to our database models:
- User type
- Product type
- Order type
- OrderItem type

### Resolvers (`resolvers.go`)
Contains all resolver functions that:
- Handle data fetching
- Process mutations
- Manage relationships between types

### Schema (`schema.go`)
Defines the GraphQL schema with:
- Root query fields for fetching users, products, orders
- Mutations for creating and updating data

### HTTP Handler (`handler.go`)
Provides an HTTP handler that:
- Processes GraphQL requests
- Executes queries against the schema
- Returns JSON responses

## 4. Running the Server

The server is configured in `api/main.go` and:
- Initializes the database connection
- Sets up the GraphQL HTTP handler with GraphiQL interface
- Starts an HTTP server on port 8081

To run the server:

```bash
cd go-graphql-ecom
go run src/api/main.go
```

The GraphQL endpoint will be available at http://localhost:8081/graphql

## 5. Testing the API

You can test the API using the provided GraphQL HTTP test file (`test/graphql.http`), which contains examples of:

- Queries for fetching users, products, and orders
- Mutations for creating users, products, orders, and order items
- Complex queries with nested relationships
- Queries and mutations with variables

If you're using VS Code with the REST Client extension or a similar tool, you can execute these requests directly from the file.

### Example Queries

Fetch all products:
```graphql
{
  products {
    id
    name
    description
    price
    inventory
  }
}
```

Create a new user:
```graphql
mutation {
  createUser(
    name: "John Doe"
    email: "john@example.com"
    password: "password123"
  ) {
    id
    name
    email
  }
}
```

Get an order with its items and related products:
```graphql
{
  order(id: 1) {
    id
    status
    total
    items {
      quantity
      price
      product {
        name
        description
      }
    }
  }
}
```

## Next Steps

Potential improvements for this project:
- Add authentication and authorization
- Implement input validation
- Add pagination for list queries
- Implement filtering and sorting
- Add error handling middleware
- Create a frontend application