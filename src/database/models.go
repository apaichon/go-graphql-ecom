package database

import (
	"database/sql"
	"errors"
)

// User represents a user in the system
type User struct {
	ID        int
	Name      string
	Email     string
	Password  string
	CreatedAt string
}

// Product represents a product in the system
type Product struct {
	ID          int
	Name        string
	Description string
	Price       float64
	Inventory   int
	CreatedAt   string
}

// Order represents an order in the system
type Order struct {
	ID        int
	UserID    int
	Status    string
	Total     float64
	CreatedAt string
	Items     []OrderItem
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID        int
	OrderID   int
	ProductID int
	Quantity  int
	Price     float64
	Product   *Product
}

// User operations

// GetUserByID retrieves a user by ID
func GetUserByID(db *sql.DB, id int) (*User, error) {
	query := `SELECT id, name, email, password, created_at FROM users WHERE id = ?`

	var user User
	err := db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetAllUsers retrieves all users
func GetAllUsers(db *sql.DB) ([]User, error) {
	query := `SELECT id, name, email, created_at FROM users`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// CreateUser creates a new user
func CreateUser(db *sql.DB, name, email, password string) (*User, error) {
	query := `INSERT INTO users (name, email, password) VALUES (?, ?, ?)`

	result, err := db.Exec(query, name, email, password)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetUserByID(db, int(id))
}

// Product operations

// GetProductByID retrieves a product by ID
func GetProductByID(db *sql.DB, id int) (*Product, error) {
	query := `SELECT id, name, description, price, inventory, created_at FROM products WHERE id = ?`

	var product Product
	err := DB.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Inventory, &product.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &product, nil
}

// GetAllProducts retrieves all products
func GetAllProducts(db *sql.DB) ([]Product, error) {
	query := `SELECT id, name, description, price, inventory, created_at FROM products`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Inventory, &product.CreatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

// CreateProduct creates a new product
func CreateProduct(db *sql.DB, name, description string, price float64, inventory int) (*Product, error) {
	query := `INSERT INTO products (name, description, price, inventory) VALUES (?, ?, ?, ?)`

	result, err := DB.Exec(query, name, description, price, inventory)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetProductByID(db, int(id))
}

// Order operations

// GetOrderByID retrieves an order by ID
func GetOrderByID(db *sql.DB, id int) (*Order, error) {
	query := `SELECT id, user_id, status, total, created_at FROM orders WHERE id = ?`

	var order Order
	err := DB.QueryRow(query, id).Scan(&order.ID, &order.UserID, &order.Status, &order.Total, &order.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	// Get order items
	items, err := GetOrderItemsByOrderID(db, order.ID)
	if err != nil {
		return nil, err
	}
	order.Items = items

	return &order, nil
}

// GetOrdersByUserID retrieves all orders for a user
func GetOrdersByUserID(db *sql.DB, userID int) ([]Order, error) {
	query := `SELECT id, user_id, status, total, created_at FROM orders WHERE user_id = ?`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(&order.ID, &order.UserID, &order.Status, &order.Total, &order.CreatedAt)
		if err != nil {
			return nil, err
		}

		// Get order items
		items, err := GetOrderItemsByOrderID(db, order.ID)
		if err != nil {
			return nil, err
		}
		order.Items = items

		orders = append(orders, order)
	}

	return orders, nil
}

// CreateOrder creates a new order
func CreateOrder(db *sql.DB, userID int, status string, total float64) (*Order, error) {
	query := `INSERT INTO orders (user_id, status, total) VALUES (?, ?, ?)`

	result, err := db.Exec(query, userID, status, total)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetOrderByID(db, int(id))
}

// OrderItem operations

// GetOrderItemsByOrderID retrieves all items for an order
func GetOrderItemsByOrderID(db *sql.DB, orderID int) ([]OrderItem, error) {
	query := `SELECT id, order_id, product_id, quantity, price FROM order_items WHERE order_id = ?`

	rows, err := db.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []OrderItem
	for rows.Next() {
		var item OrderItem
		err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price)
		if err != nil {
			return nil, err
		}

		// Get associated product
		product, err := GetProductByID(db, item.ProductID)
		if err != nil {
			return nil, err
		}
		item.Product = product

		items = append(items, item)
	}

	return items, nil
}

// AddOrderItem adds an item to an order
func AddOrderItem(db *sql.DB, orderID, productID, quantity int, price float64) (*OrderItem, error) {
	query := `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)`

	result, err := DB.Exec(query, orderID, productID, quantity, price)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Retrieve the created order item
	var orderItem OrderItem
	err = DB.QueryRow("SELECT id, order_id, product_id, quantity, price FROM order_items WHERE id = ?", id).
		Scan(&orderItem.ID, &orderItem.OrderID, &orderItem.ProductID, &orderItem.Quantity, &orderItem.Price)
	if err != nil {
		return nil, err
	}

	// Get associated product
	product, err := GetProductByID(db, orderItem.ProductID)
	if err != nil {
		return nil, err
	}
	orderItem.Product = product

	return &orderItem, nil
}
