package main

import (
	"fmt"
	"log"
	"net/http"

	"go-graphql-ecom/database"
	"go-graphql-ecom/graphql"

	"github.com/graphql-go/handler"
)

func main() {
	// Initialize database
	err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// Create a GraphQL HTTP handler with our schema
	h := handler.New(&handler.Config{
		Schema:   &graphql.Schema,
		Pretty:   true,
		GraphiQL: true,
	})

	// Set up GraphQL endpoint
	http.Handle("/graphql", h)

	// Start server
	fmt.Println("Server is running on http://localhost:8081/graphql")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
