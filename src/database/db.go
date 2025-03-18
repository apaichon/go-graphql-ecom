package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// DB is the database connection
var (
	DB   *sql.DB
	once sync.Once
)

// InitDB initializes the database connection
func InitDB() error {
	var err error

	once.Do(func() {
		// Check if the database file exists, and if not, create it
		if _, err := os.Stat("../../data/ecommerce.db"); os.IsNotExist(err) {
			file, err := os.Create("../../data/ecommerce.db")
			if err != nil {
				log.Fatal(err)
			}
			file.Close()
		}

		DB, err = sql.Open("sqlite3", "../../data/ecommerce.db")
		if err != nil {
			log.Fatal(err)
		}

		// Test the connection
		err = DB.Ping()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Connected to SQLite database")

		// Create tables if they don't exist
		createTables()
	})

	return err
}

// GetDB returns the database instance, initializing it if necessary
func GetDB() *sql.DB {
	if DB == nil {
		err := InitDB()
		if err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}
	}
	return DB
}

// createTables creates all necessary tables for our e-commerce application
func createTables() {
	// Create Users table
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := DB.Exec(usersTable)
	if err != nil {
		log.Fatal(err)
	}

	// Create Products table
	productsTable := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		price REAL NOT NULL,
		inventory INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = DB.Exec(productsTable)
	if err != nil {
		log.Fatal(err)
	}

	// Create Orders table
	ordersTable := `
	CREATE TABLE IF NOT EXISTS orders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		status TEXT NOT NULL,
		total REAL NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users (id)
	);
	`
	_, err = DB.Exec(ordersTable)
	if err != nil {
		log.Fatal(err)
	}

	// Create OrderItems table
	orderItemsTable := `
	CREATE TABLE IF NOT EXISTS order_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		order_id INTEGER NOT NULL,
		product_id INTEGER NOT NULL,
		quantity INTEGER NOT NULL,
		price REAL NOT NULL,
		FOREIGN KEY (order_id) REFERENCES orders (id),
		FOREIGN KEY (product_id) REFERENCES products (id)
	);
	`
	_, err = DB.Exec(orderItemsTable)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Tables created successfully")
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
		DB = nil
	}
}

// GetAllOrders retrieves all orders from the database
func GetAllOrders() ([]*Order, error) {
	rows, err := DB.Query("SELECT id, user_id, status, total, created_at FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*Order
	for rows.Next() {
		order := &Order{}
		err := rows.Scan(&order.ID, &order.UserID, &order.Status, &order.Total, &order.CreatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}
