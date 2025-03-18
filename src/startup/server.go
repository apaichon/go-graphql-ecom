package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"go-graphql-ecom/database"
)

// Define GraphQL types
var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id":    &graphql.Field{Type: graphql.Int},
		"name":  &graphql.Field{Type: graphql.String},
		"email": &graphql.Field{Type: graphql.String},
	},
})

var productType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Product",
	Fields: graphql.Fields{
		"id":          &graphql.Field{Type: graphql.Int},
		"name":        &graphql.Field{Type: graphql.String},
		"description": &graphql.Field{Type: graphql.String},
		"price":       &graphql.Field{Type: graphql.Float},
		"inventory":   &graphql.Field{Type: graphql.Int},
	},
})

// Mock data
var users = []map[string]interface{}{
	{"id": 1, "name": "John Doe", "email": "john@example.com"},
	{"id": 2, "name": "Jane Smith", "email": "jane@example.com"},
}

var products = []map[string]interface{}{
	{"id": 1, "name": "Laptop", "description": "High-performance laptop", "price": 999.99, "inventory": 10},
	{"id": 2, "name": "Smartphone", "description": "Latest smartphone", "price": 699.99, "inventory": 20},
}

// InitSchema initializes the GraphQL schema
func InitSchema() (*graphql.Schema, error) {
	// Define root query
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(int)
					if !ok {
						return nil, nil
					}

					/*for _, user := range users {
						if user["id"] == id {
							return user, nil
						}
					}*/
					db := database.GetDB()
					user, err := database.GetUserByID(db,id)
					if err != nil {
						return nil, err
					}
					return user, nil
				},
			},
			"users": &graphql.Field{
				Type: graphql.NewList(userType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return users, nil
				},
			},
			"product": &graphql.Field{
				Type: productType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(int)
					if !ok {
						return nil, nil
					}

					for _, product := range products {
						if product["id"] == id {
							return product, nil
						}
					}
					return nil, nil
				},
			},
			"products": &graphql.Field{
				Type: graphql.NewList(productType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return products, nil
				},
			},
		},
	})

	// Create schema config
	schemaConfig := graphql.SchemaConfig{
		Query: rootQuery,
	}

	// Create schema
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		return nil, err
	}

	return &schema, nil
}

// ExecuteQuery processes a GraphQL query and returns the result
func ExecuteQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})

	if len(result.Errors) > 0 {
		fmt.Printf("GraphQL query errors: %v\n", result.Errors)
	}

	return result
}

// GraphQLHandler handles HTTP requests for GraphQL queries
func GraphQLHandler(schema *graphql.Schema) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for browser clients
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Only accept POST requests
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Method not allowed: %s", r.Method)
			return
		}

		// Parse the request body
		var requestBody struct {
			Query     string                 `json:"query"`
			Variables map[string]interface{} `json:"variables"`
		}

		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error parsing request body: %v", err)
			return
		}

		// Execute the GraphQL query
		result := graphql.Do(graphql.Params{
			Schema:         *schema,
			RequestString:  requestBody.Query,
			VariableValues: requestBody.Variables,
		})

		// Return the result as JSON
		w.Header().Set("Content-Type", "application/json")
		if len(result.Errors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
		}

		json.NewEncoder(w).Encode(result)
	}
}

// PlaygroundHandler serves the GraphQL Playground interface
func PlaygroundHandler(endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		playground := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>GraphQL Playground</title>
  <meta name="viewport" content="user-scalable=no, initial-scale=1.0, minimum-scale=1.0, maximum-scale=1.0, minimal-ui">
  <link rel="stylesheet" href="//cdn.jsdelivr.net/npm/graphql-playground-react/build/static/css/index.css" />
  <link rel="shortcut icon" href="//cdn.jsdelivr.net/npm/graphql-playground-react/build/favicon.png" />
  <script src="//cdn.jsdelivr.net/npm/graphql-playground-react/build/static/js/middleware.js"></script>
</head>
<body>
  <div id="root">
    <style>
      body {
        background-color: rgb(23, 42, 58);
        font-family: Open Sans, sans-serif;
        height: 90vh;
      }
      #root {
        height: 100vh;
        width: 100vw;
        display: flex;
        align-items: center;
        justify-content: center;
      }
      .loading {
        font-size: 32px;
        font-weight: 200;
        color: rgba(255, 255, 255, .6);
        margin-left: 20px;
      }
      img {
        width: 78px;
        height: 78px;
      }
      .title {
        font-weight: 400;
      }
    </style>
    <img src='//cdn.jsdelivr.net/npm/graphql-playground-react/build/logo.png' alt=''>
    <div class="loading"> Loading
      <span class="title">GraphQL Playground</span>
    </div>
  </div>
  <script>window.addEventListener('load', function (event) {
      GraphQLPlayground.init(document.getElementById('root'), {
        endpoint: '%s'
      })
    })</script>
</body>
</html>
`, endpoint)
		fmt.Fprint(w, playground)
	}
}

func main() {
	// Initialize GraphQL schema
	schema, err := InitSchema()
	if err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	// Create a router
	http.HandleFunc("/graphql", GraphQLHandler(schema))
	http.HandleFunc("/playground", PlaygroundHandler("/graphql"))

	// Start the server
	port := ":4000"
	fmt.Printf("Server running on http://localhost%s\n", port)
	fmt.Printf("GraphQL endpoint: http://localhost%s/graphql\n", port)
	fmt.Printf("GraphQL Playground: http://localhost%s/playground\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
