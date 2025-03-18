package graphql

import (
	"errors"
	"fmt"

	"go-graphql-ecom/database"

	"github.com/graphql-go/graphql"
)

// User resolvers
func getUserResolver(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(int)
	if !ok {
		return nil, errors.New("invalid user ID")
	}
	return database.GetUserByID(database.GetDB(), id)
}

func getAllUsersResolver(p graphql.ResolveParams) (interface{}, error) {
	return database.GetAllUsers(database.GetDB())
}

func createUserResolver(p graphql.ResolveParams) (interface{}, error) {
	name := p.Args["name"].(string)
	email := p.Args["email"].(string)
	password := p.Args["password"].(string)

	return database.CreateUser(database.GetDB(), name, email, password)
}

// Product resolvers
func getProductResolver(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(int)
	if !ok {
		return nil, errors.New("invalid product ID")
	}
	return database.GetProductByID(database.GetDB(), id)
}

func getAllProductsResolver(p graphql.ResolveParams) (interface{}, error) {
	return database.GetAllProducts(database.GetDB())
}

func createProductResolver(p graphql.ResolveParams) (interface{}, error) {
	name := p.Args["name"].(string)
	description, _ := p.Args["description"].(string)
	price := p.Args["price"].(float64)
	inventory := p.Args["inventory"].(int)

	return database.CreateProduct(database.GetDB(), name, description, price, inventory)
}

// Order resolvers
func getOrderResolver(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(int)
	if !ok {
		return nil, errors.New("invalid order ID")
	}
	return database.GetOrderByID(database.GetDB(), id)
}

func getAllOrdersResolver(p graphql.ResolveParams) (interface{}, error) {
	return database.GetAllOrders()
}

func createOrderResolver(p graphql.ResolveParams) (interface{}, error) {
	userID := p.Args["user_id"].(int)
	status := p.Args["status"].(string)
	total := p.Args["total"].(float64)

	return database.CreateOrder(database.GetDB(), userID, status, total)
}

func updateOrderStatusResolver(p graphql.ResolveParams) (interface{}, error) {
	id := p.Args["id"].(int)
	status := p.Args["status"].(string)

	db := database.GetDB()
	_, err := db.Exec("UPDATE orders SET status = ? WHERE id = ?", status, id)
	if err != nil {
		return nil, err
	}

	return database.GetOrderByID(db, id)
}

// OrderItem resolvers
func addOrderItemResolver(p graphql.ResolveParams) (interface{}, error) {
	orderID := p.Args["order_id"].(int)
	productID := p.Args["product_id"].(int)
	quantity := p.Args["quantity"].(int)
	price := p.Args["price"].(float64)

	// Check if product has enough inventory
	db := database.GetDB()
	product, err := database.GetProductByID(db, productID)
	if err != nil {
		return nil, err
	}

	if product.Inventory < quantity {
		return nil, fmt.Errorf("not enough inventory for product %d", productID)
	}

	return database.AddOrderItem(db, orderID, productID, quantity, price)
}

// Relationship resolvers
func getProductFromOrderItemResolver(p graphql.ResolveParams) (interface{}, error) {
	if orderItem, ok := p.Source.(*database.OrderItem); ok {
		return orderItem.Product, nil
	}
	return nil, errors.New("failed to get product from order item")
}

func getItemsFromOrderResolver(p graphql.ResolveParams) (interface{}, error) {
	if order, ok := p.Source.(*database.Order); ok {
		return order.Items, nil
	}
	return nil, errors.New("failed to get items from order")
}
