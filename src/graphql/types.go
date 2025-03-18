package graphql

import (

	"github.com/graphql-go/graphql"
)


// Define GraphQL types
var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"email": &graphql.Field{
			Type: graphql.String,
		},
		"created_at": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var productType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Product",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
		"price": &graphql.Field{
			Type: graphql.Float,
		},
		"inventory": &graphql.Field{
			Type: graphql.Int,
		},
		"created_at": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var orderItemType = graphql.NewObject(graphql.ObjectConfig{
	Name: "OrderItem",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"order_id": &graphql.Field{
			Type: graphql.Int,
		},
		"product_id": &graphql.Field{
			Type: graphql.Int,
		},
		"quantity": &graphql.Field{
			Type: graphql.Int,
		},
		"price": &graphql.Field{
			Type: graphql.Float,
		},
		"product": &graphql.Field{
			Type: productType,
			Resolve: getProductFromOrderItemResolver,
		},
	},
})

var orderType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Order",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"user_id": &graphql.Field{
			Type: graphql.Int,
		},
		"status": &graphql.Field{
			Type: graphql.String,
		},
		"total": &graphql.Field{
			Type: graphql.Float,
		},
		"created_at": &graphql.Field{
			Type: graphql.String,
		},
		"items": &graphql.Field{
			Type: graphql.NewList(orderItemType),
			Resolve: getItemsFromOrderResolver,
		},
	},
})
