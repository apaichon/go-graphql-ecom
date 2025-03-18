package graphql

import (
	"encoding/json"
	"net/http"

	"github.com/graphql-go/graphql"

)

type postData struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

// Handler handles GraphQL HTTP requests
func Handler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var data postData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": []map[string]interface{}{
				{
					"message": "Invalid request body",
				},
			},
		})
		return
	}

	// Execute the GraphQL query
	result := graphql.Do(graphql.Params{
		Schema:         Schema,
		RequestString:  data.Query,
		VariableValues: data.Variables,
	})

	// Set content type and return the result
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
