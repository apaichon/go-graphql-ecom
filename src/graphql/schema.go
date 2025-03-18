package graphql

import (
	"github.com/graphql-go/graphql"
)

// Define root query
var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"user": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: getUserResolver,
		},
		"users": &graphql.Field{
			Type:    graphql.NewList(userType),
			Resolve: getAllUsersResolver,
		},
		"product": &graphql.Field{
			Type: productType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: getProductResolver,
		},
		"products": &graphql.Field{
			Type:    graphql.NewList(productType),
			Resolve: getAllProductsResolver,
		},
		"order": &graphql.Field{
			Type: orderType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: getOrderResolver,
		},
		"orders": &graphql.Field{
			Type:    graphql.NewList(orderType),
			Resolve: getAllOrdersResolver,
		},
	},
})

// Define mutations
var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"createUser": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"email": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"password": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: createUserResolver,
		},
		"createProduct": &graphql.Field{
			Type: productType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"description": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"price": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Float),
				},
				"inventory": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: createProductResolver,
		},
		"createOrder": &graphql.Field{
			Type: orderType,
			Args: graphql.FieldConfigArgument{
				"user_id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"status": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"total": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Float),
				},
			},
			Resolve: createOrderResolver,
		},
		"addOrderItem": &graphql.Field{
			Type: orderItemType,
			Args: graphql.FieldConfigArgument{
				"order_id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"product_id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"quantity": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"price": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Float),
				},
			},
			Resolve: addOrderItemResolver,
		},
		"updateOrderStatus": &graphql.Field{
			Type: orderType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"status": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: updateOrderStatusResolver,
		},
	},
})

// Create schema
var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})
