# pgxutil

[![Build Status](https://github.com/allisson/pgxutil/workflows/Release/badge.svg)](https://github.com/allisson/pgxutil/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/allisson/pgxutil/v2)](https://goreportcard.com/report/github.com/allisson/pgxutil/v2)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/allisson/pgxutil/v2)

High-level helper functions for PostgreSQL database operations using the [pgx](https://github.com/jackc/pgx) driver. This library wraps [sqlquery](https://github.com/allisson/sqlquery) for SQL generation and [scany](https://github.com/georgysavva/scany) for row scanning, providing type-safe and convenient methods for common CRUD operations.

## Features

- **Type-safe CRUD operations** - Insert, Select, Update, and Delete with struct mapping
- **Flexible query building** - Rich filtering, sorting, pagination, and field selection
- **Works with any pgx type** - Supports `*pgxpool.Pool`, `*pgx.Conn`, and `pgx.Tx`
- **Minimal boilerplate** - Reduces repetitive SQL query writing
- **PostgreSQL-focused** - Takes advantage of PostgreSQL-specific features

## Installation

```bash
go get github.com/allisson/pgxutil/v2
```

## Table of Contents

- [Quick Start](#quick-start)
- [Complete Example](#complete-example)
- [Insert Operations](#insert-operations)
- [Select Operations](#select-operations)
- [Update Operations](#update-operations)
- [Delete Operations](#delete-operations)
- [Advanced Filtering](#advanced-filtering)
- [Pagination and Sorting](#pagination-and-sorting)
- [Working with Transactions](#working-with-transactions)
- [Working with Connection Pools](#working-with-connection-pools)
- [Using FOR UPDATE](#using-for-update)
- [Field Selection](#field-selection)

## Quick Start

```go
package main

import (
	"context"
	"log"

	"github.com/allisson/pgxutil/v2"
	"github.com/jackc/pgx/v5"
)

type User struct {
	ID    int    `db:"id"`
	Email string `db:"email" fieldtag:"insert,update"`
	Name  string `db:"name" fieldtag:"insert,update"`
}

func main() {
	ctx := context.Background()
	conn, _ := pgx.Connect(ctx, "postgres://user:pass@localhost/mydb?sslmode=disable")
	defer conn.Close(ctx)

	// Insert a user
	user := User{Email: "john@example.com", Name: "John Doe"}
	if err := pgxutil.Insert(ctx, conn, "insert", "users", &user); err != nil {
		log.Fatal(err)
	}

	// Get a user by email
	options := pgxutil.NewFindOptions().WithFilter("email", "john@example.com")
	if err := pgxutil.Get(ctx, conn, "users", options, &user); err != nil {
		log.Fatal(err)
	}

	log.Printf("User: %+v", user)
}
```

## Complete Example

Here's a comprehensive example showing all major operations:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/allisson/pgxutil/v2"
	"github.com/jackc/pgx/v5"
)

type Player struct {
	ID   int    `db:"id"`
	Name string `db:"name" fieldtag:"insert,update"`
	Age  int    `db:"age" fieldtag:"insert,update"`
	Team string `db:"team" fieldtag:"insert,update"`
}

func main() {
	ctx := context.Background()
	
	// Connect to database
	// Run a database with docker: docker run --name pgxutil-test --restart unless-stopped -e POSTGRES_USER=user -e POSTGRES_PASSWORD=password -e POSTGRES_DB=pgxutil -p 5432:5432 -d postgres:16-alpine
	conn, err := pgx.Connect(ctx, "postgres://user:password@localhost/pgxutil?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	// Create table
	_, err = conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS players(
			id SERIAL PRIMARY KEY,
			name VARCHAR NOT NULL,
			age INTEGER NOT NULL,
			team VARCHAR NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Insert players
	r9 := Player{Name: "Ronaldo Fenômeno", Age: 44, Team: "Real Madrid"}
	r10 := Player{Name: "Ronaldinho Gaúcho", Age: 41, Team: "Barcelona"}
	messi := Player{Name: "Lionel Messi", Age: 36, Team: "Inter Miami"}
	
	tag := "insert" // uses fields with fieldtag:"insert"
	if err := pgxutil.Insert(ctx, conn, tag, "players", &r9); err != nil {
		log.Fatal(err)
	}
	if err := pgxutil.Insert(ctx, conn, tag, "players", &r10); err != nil {
		log.Fatal(err)
	}
	if err := pgxutil.Insert(ctx, conn, tag, "players", &messi); err != nil {
		log.Fatal(err)
	}

	// Get single player by name
	findOptions := pgxutil.NewFindOptions().WithFilter("name", "Lionel Messi")
	if err := pgxutil.Get(ctx, conn, "players", findOptions, &messi); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found player: %+v\n", messi)

	// Select all players ordered by name
	players := []*Player{}
	findAllOptions := pgxutil.NewFindAllOptions().
		WithLimit(10).
		WithOffset(0).
		WithOrderBy("name asc")
	if err := pgxutil.Select(ctx, conn, "players", findAllOptions, &players); err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("\nAll players:")
	for _, p := range players {
		fmt.Printf("  %+v\n", p)
	}

	// Update player
	tag = "update" // uses fields with fieldtag:"update"
	messi.Team = "Argentina National Team"
	if err := pgxutil.Update(ctx, conn, tag, "players", messi.ID, &messi); err != nil {
		log.Fatal(err)
	}

	// Delete player
	if err := pgxutil.Delete(ctx, conn, "players", r9.ID); err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nOperations completed successfully!")
}
```

## Insert Operations

### Basic Insert

```go
type Product struct {
	ID    int     `db:"id"`
	Name  string  `db:"name" fieldtag:"insert"`
	Price float64 `db:"price" fieldtag:"insert"`
	Stock int     `db:"stock" fieldtag:"insert"`
}

product := Product{
	Name:  "Laptop",
	Price: 999.99,
	Stock: 10,
}

err := pgxutil.Insert(ctx, conn, "insert", "products", &product)
```

### Insert with Selective Fields

Use different field tags to control which fields get inserted:

```go
type User struct {
	ID        int       `db:"id"`
	Email     string    `db:"email" fieldtag:"insert,register"`
	Password  string    `db:"password" fieldtag:"insert,register"`
	Name      string    `db:"name" fieldtag:"insert,update,register"`
	IsActive  bool      `db:"is_active" fieldtag:"insert,update"`
	CreatedAt time.Time `db:"created_at"`
}

// During user registration (only email, password, name)
user := User{Email: "user@example.com", Password: "hashed", Name: "Jane"}
err := pgxutil.Insert(ctx, conn, "register", "users", &user)

// During admin creation (includes is_active)
adminUser := User{Email: "admin@example.com", Password: "hashed", Name: "Admin", IsActive: true}
err = pgxutil.Insert(ctx, conn, "insert", "users", &adminUser)
```

### Batch Insert

```go
orders := []Order{
	{CustomerID: 1, TotalAmount: 100.50, Status: "pending"},
	{CustomerID: 2, TotalAmount: 250.00, Status: "pending"},
	{CustomerID: 3, TotalAmount: 75.25, Status: "pending"},
}

for _, order := range orders {
	if err := pgxutil.Insert(ctx, conn, "insert", "orders", &order); err != nil {
		log.Printf("Failed to insert order: %v", err)
	}
}
```

## Select Operations

### Get Single Record

```go
type Book struct {
	ID     int    `db:"id"`
	Title  string `db:"title"`
	Author string `db:"author"`
	ISBN   string `db:"isbn"`
}

// Find by ID
options := pgxutil.NewFindOptions().WithFilter("id", 42)
var book Book
err := pgxutil.Get(ctx, conn, "books", options, &book)

// Find by ISBN
options = pgxutil.NewFindOptions().WithFilter("isbn", "978-0-123456-78-9")
err = pgxutil.Get(ctx, conn, "books", options, &book)

// Find by multiple conditions
options = pgxutil.NewFindOptions().
	WithFilter("author", "J.K. Rowling").
	WithFilter("title.like", "%Harry Potter%")
err = pgxutil.Get(ctx, conn, "books", options, &book)
```

### Select Multiple Records

```go
// Get all books
var books []Book
options := pgxutil.NewFindAllOptions()
err := pgxutil.Select(ctx, conn, "books", options, &books)

// Get books with filtering
options = pgxutil.NewFindAllOptions().
	WithFilter("author", "J.K. Rowling").
	WithOrderBy("title asc")
err = pgxutil.Select(ctx, conn, "books", options, &books)

// Get recent books
options = pgxutil.NewFindAllOptions().
	WithFilter("published_year.gte", 2020).
	WithLimit(10).
	WithOrderBy("published_year desc")
err = pgxutil.Select(ctx, conn, "books", options, &books)
```

### Select with Specific Fields

```go
type BookSummary struct {
	Title  string `db:"title"`
	Author string `db:"author"`
}

var summaries []BookSummary
options := pgxutil.NewFindAllOptions().
	WithFields([]string{"title", "author"}).
	WithOrderBy("title asc")
err := pgxutil.Select(ctx, conn, "books", options, &summaries)
```

## Update Operations

### Update by ID

```go
type Article struct {
	ID        int    `db:"id"`
	Title     string `db:"title" fieldtag:"update"`
	Content   string `db:"content" fieldtag:"update"`
	Published bool   `db:"published" fieldtag:"update"`
}

// Fetch the article
options := pgxutil.NewFindOptions().WithFilter("id", 10)
var article Article
err := pgxutil.Get(ctx, conn, "articles", options, &article)

// Modify and update
article.Published = true
article.Title = "Updated Title"
err = pgxutil.Update(ctx, conn, "update", "articles", article.ID, &article)
```

### Update with Options (Advanced)

Update multiple records with complex conditions:

```go
// Update all unpublished articles older than 30 days to archived status
updateOpts := pgxutil.NewUpdateOptions().
	WithSet("status", "archived").
	WithSet("archived_at", time.Now()).
	WithFilter("published", false).
	WithFilter("created_at.lt", time.Now().AddDate(0, 0, -30))

err := pgxutil.UpdateWithOptions(ctx, conn, "articles", updateOpts)
```

### Conditional Updates

```go
// Mark orders as shipped if they're in processing status
updateOpts := pgxutil.NewUpdateOptions().
	WithSet("status", "shipped").
	WithSet("shipped_at", time.Now()).
	WithFilter("status", "processing").
	WithFilter("payment_status", "paid")

err := pgxutil.UpdateWithOptions(ctx, conn, "orders", updateOpts)
```

### Update with Complex Filters

```go
// Increase price for products low in stock
updateOpts := pgxutil.NewUpdateOptions().
	WithSet("price", "price * 1.1"). // Use SQL expression
	WithFilter("stock.lt", 5).
	WithFilter("category", "electronics")

err := pgxutil.UpdateWithOptions(ctx, conn, "products", updateOpts)
```

## Delete Operations

### Delete by ID

```go
// Simple delete
err := pgxutil.Delete(ctx, conn, "users", 123)

// Delete and check error
if err := pgxutil.Delete(ctx, conn, "sessions", sessionID); err != nil {
	log.Printf("Failed to delete session: %v", err)
}
```

### Delete with Options (Advanced)

Delete multiple records with conditions:

```go
// Delete old logs
deleteOpts := pgxutil.NewDeleteOptions().
	WithFilter("created_at.lt", time.Now().AddDate(0, 0, -90))

err := pgxutil.DeleteWithOptions(ctx, conn, "logs", deleteOpts)
```

### Conditional Deletes

```go
// Delete expired sessions
deleteOpts := pgxutil.NewDeleteOptions().
	WithFilter("expires_at.lt", time.Now())

err := pgxutil.DeleteWithOptions(ctx, conn, "sessions", deleteOpts)

// Delete inactive users who never verified email
deleteOpts = pgxutil.NewDeleteOptions().
	WithFilter("is_active", false).
	WithFilter("email_verified", false).
	WithFilter("created_at.lt", time.Now().AddDate(0, -6, 0)) // 6 months ago

err = pgxutil.DeleteWithOptions(ctx, conn, "users", deleteOpts)
```

### Bulk Delete

```go
// Delete specific IDs
deleteOpts := pgxutil.NewDeleteOptions().
	WithFilter("id.in", "1,2,3,4,5")

err := pgxutil.DeleteWithOptions(ctx, conn, "temp_records", deleteOpts)
```

## Advanced Filtering

pgxutil supports a wide range of filter operators:

### Basic Filters

```go
options := pgxutil.NewFindAllOptions().
	WithFilter("status", "active").        // WHERE status = 'active'
	WithFilter("price", 99.99).            // WHERE price = 99.99
	WithFilter("deleted_at", nil)          // WHERE deleted_at IS NULL
```

### Comparison Operators

```go
options := pgxutil.NewFindAllOptions().
	WithFilter("age.gt", 18).              // WHERE age > 18
	WithFilter("age.gte", 18).             // WHERE age >= 18
	WithFilter("age.lt", 65).              // WHERE age < 65
	WithFilter("age.lte", 65).             // WHERE age <= 65
	WithFilter("status.not", "banned")     // WHERE status <> 'banned'
```

### IN and NOT IN

```go
options := pgxutil.NewFindAllOptions().
	WithFilter("id.in", "1,2,3,4,5").      // WHERE id IN (1,2,3,4,5)
	WithFilter("status.notin", "deleted,banned") // WHERE status NOT IN ('deleted','banned')
```

### LIKE Pattern Matching

```go
options := pgxutil.NewFindAllOptions().
	WithFilter("email.like", "%@gmail.com"). // WHERE email LIKE '%@gmail.com'
	WithFilter("name.like", "John%")         // WHERE name LIKE 'John%'
```

### NULL Checks

```go
options := pgxutil.NewFindAllOptions().
	WithFilter("deleted_at.null", true).   // WHERE deleted_at IS NULL
	WithFilter("verified_at.null", false)  // WHERE verified_at IS NOT NULL
```

### Combining Multiple Filters

```go
// Complex query: Active users from specific countries, aged 18-65, email verified
options := pgxutil.NewFindAllOptions().
	WithFilter("is_active", true).
	WithFilter("country.in", "US,CA,UK").
	WithFilter("age.gte", 18).
	WithFilter("age.lte", 65).
	WithFilter("email_verified_at.null", false).
	WithOrderBy("created_at desc").
	WithLimit(50)

var users []User
err := pgxutil.Select(ctx, conn, "users", options, &users)
```

## Pagination and Sorting

### Basic Pagination

```go
page := 1
pageSize := 20
offset := (page - 1) * pageSize

options := pgxutil.NewFindAllOptions().
	WithLimit(pageSize).
	WithOffset(offset).
	WithOrderBy("created_at desc")

var items []Item
err := pgxutil.Select(ctx, conn, "items", options, &items)
```

### Pagination Helper Function

```go
func GetPage(ctx context.Context, db pgxutil.Querier, page, pageSize int) ([]Product, error) {
	offset := (page - 1) * pageSize
	
	options := pgxutil.NewFindAllOptions().
		WithLimit(pageSize).
		WithOffset(offset).
		WithOrderBy("name asc")
	
	var products []Product
	err := pgxutil.Select(ctx, db, "products", options, &products)
	return products, err
}

// Usage
products, err := GetPage(ctx, conn, 1, 25) // First page, 25 items per page
```

### Multiple Sort Fields

```go
options := pgxutil.NewFindAllOptions().
	WithOrderBy("priority desc, created_at asc")

var tasks []Task
err := pgxutil.Select(ctx, conn, "tasks", options, &tasks)
```

### Cursor-based Pagination

```go
// Get next page after last seen ID
lastSeenID := 100

options := pgxutil.NewFindAllOptions().
	WithFilter("id.gt", lastSeenID).
	WithOrderBy("id asc").
	WithLimit(20)

var posts []Post
err := pgxutil.Select(ctx, conn, "posts", options, &posts)
```

## Working with Transactions

pgxutil works seamlessly with pgx transactions:

```go
func TransferFunds(ctx context.Context, pool *pgxpool.Pool, fromAccountID, toAccountID int, amount float64) error {
	// Begin transaction
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Deduct from source account
	updateOpts := pgxutil.NewUpdateOptions().
		WithSet("balance", fmt.Sprintf("balance - %f", amount)).
		WithFilter("id", fromAccountID).
		WithFilter("balance.gte", amount) // Ensure sufficient funds

	if err := pgxutil.UpdateWithOptions(ctx, tx, "accounts", updateOpts); err != nil {
		return fmt.Errorf("failed to deduct from source account: %w", err)
	}

	// Add to destination account
	updateOpts = pgxutil.NewUpdateOptions().
		WithSet("balance", fmt.Sprintf("balance + %f", amount)).
		WithFilter("id", toAccountID)

	if err := pgxutil.UpdateWithOptions(ctx, tx, "accounts", updateOpts); err != nil {
		return fmt.Errorf("failed to add to destination account: %w", err)
	}

	// Record transaction
	transaction := Transaction{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        amount,
		CreatedAt:     time.Now(),
	}
	if err := pgxutil.Insert(ctx, tx, "insert", "transactions", &transaction); err != nil {
		return fmt.Errorf("failed to record transaction: %w", err)
	}

	// Commit transaction
	return tx.Commit(ctx)
}
```

### Transaction with SELECT FOR UPDATE

```go
func ProcessOrder(ctx context.Context, pool *pgxpool.Pool, orderID int) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Lock the order for processing
	options := pgxutil.NewFindOptions().
		WithFilter("id", orderID).
		WithForUpdate("NOWAIT") // Fail immediately if locked by another transaction

	var order Order
	if err := pgxutil.Get(ctx, tx, "orders", options, &order); err != nil {
		return fmt.Errorf("failed to lock order: %w", err)
	}

	// Check if already processed
	if order.Status == "completed" {
		return fmt.Errorf("order already processed")
	}

	// Process the order...
	order.Status = "completed"
	order.ProcessedAt = time.Now()

	if err := pgxutil.Update(ctx, tx, "update", "orders", order.ID, &order); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
```

## Working with Connection Pools

Using pgxpool for production applications:

```go
package main

import (
	"context"
	"log"

	"github.com/allisson/pgxutil/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	// Create connection pool
	config, err := pgxpool.ParseConfig("postgres://user:pass@localhost/mydb?pool_max_conns=10")
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Use pool with pgxutil
	options := pgxutil.NewFindAllOptions().WithLimit(100)
	var users []User
	if err := pgxutil.Select(ctx, pool, "users", options, &users); err != nil {
		log.Fatal(err)
	}

	log.Printf("Found %d users", len(users))
}
```

### Repository Pattern with Pool

```go
type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, user *User) error {
	return pgxutil.Insert(ctx, r.pool, "insert", "users", user)
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*User, error) {
	options := pgxutil.NewFindOptions().WithFilter("id", id)
	var user User
	err := pgxutil.Get(ctx, r.pool, "users", options, &user)
	return &user, err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	options := pgxutil.NewFindOptions().WithFilter("email", email)
	var user User
	err := pgxutil.Get(ctx, r.pool, "users", options, &user)
	return &user, err
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]User, error) {
	options := pgxutil.NewFindAllOptions().
		WithLimit(limit).
		WithOffset(offset).
		WithOrderBy("created_at desc")
	
	var users []User
	err := pgxutil.Select(ctx, r.pool, "users", options, &users)
	return users, err
}

func (r *UserRepository) Update(ctx context.Context, user *User) error {
	return pgxutil.Update(ctx, r.pool, "update", "users", user.ID, user)
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	return pgxutil.Delete(ctx, r.pool, "users", id)
}
```

## Using FOR UPDATE

Lock rows for update in concurrent scenarios:

### FOR UPDATE

```go
// Lock row for update (blocks until lock is available)
options := pgxutil.NewFindOptions().
	WithFilter("id", 1).
	WithForUpdate("")

var account Account
err := pgxutil.Get(ctx, tx, "accounts", options, &account)
```

### FOR UPDATE NOWAIT

```go
// Fail immediately if row is locked
options := pgxutil.NewFindOptions().
	WithFilter("id", 1).
	WithForUpdate("NOWAIT")

var account Account
err := pgxutil.Get(ctx, tx, "accounts", options, &account)
if err != nil {
	// Handle lock contention
}
```

### FOR UPDATE SKIP LOCKED

```go
// Skip locked rows (useful for queue processing)
options := pgxutil.NewFindAllOptions().
	WithFilter("status", "pending").
	WithLimit(10).
	WithForUpdate("SKIP LOCKED")

var jobs []Job
err := pgxutil.Select(ctx, tx, "jobs", options, &jobs)
// Process only unlocked jobs
```

### Queue Worker Example

```go
func ProcessNextJob(ctx context.Context, pool *pgxpool.Pool) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Get next available job (skip locked ones)
	options := pgxutil.NewFindAllOptions().
		WithFilter("status", "pending").
		WithOrderBy("created_at asc").
		WithLimit(1).
		WithForUpdate("SKIP LOCKED")

	var jobs []Job
	if err := pgxutil.Select(ctx, tx, "jobs", options, &jobs); err != nil {
		return err
	}

	if len(jobs) == 0 {
		return nil // No jobs available
	}

	job := jobs[0]

	// Mark as processing
	job.Status = "processing"
	job.StartedAt = time.Now()
	if err := pgxutil.Update(ctx, tx, "update", "jobs", job.ID, &job); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	// Process job outside transaction...
	return nil
}
```

## Field Selection

Select only specific fields to improve performance:

```go
// Select only ID and email
type UserBasic struct {
	ID    int    `db:"id"`
	Email string `db:"email"`
}

options := pgxutil.NewFindAllOptions().
	WithFields([]string{"id", "email"}).
	WithFilter("is_active", true)

var users []UserBasic
err := pgxutil.Select(ctx, conn, "users", options, &users)
```

### Projection for API Responses

```go
type ProductListItem struct {
	ID    int     `db:"id"`
	Name  string  `db:"name"`
	Price float64 `db:"price"`
}

func GetProductList(ctx context.Context, db pgxutil.Querier) ([]ProductListItem, error) {
	options := pgxutil.NewFindAllOptions().
		WithFields([]string{"id", "name", "price"}). // Skip description, images, etc.
		WithFilter("is_available", true).
		WithOrderBy("name asc")
	
	var products []ProductListItem
	err := pgxutil.Select(ctx, db, "products", options, &products)
	return products, err
}
```

## Complete Filter Reference

All available filter operators:

| Operator | Example | SQL Generated |
|----------|---------|---------------|
| (none) | `WithFilter("status", "active")` | `WHERE status = 'active'` |
| `nil` | `WithFilter("deleted_at", nil)` | `WHERE deleted_at IS NULL` |
| `.in` | `WithFilter("id.in", "1,2,3")` | `WHERE id IN (1,2,3)` |
| `.notin` | `WithFilter("status.notin", "a,b")` | `WHERE status NOT IN ('a','b')` |
| `.not` | `WithFilter("status.not", "banned")` | `WHERE status <> 'banned'` |
| `.gt` | `WithFilter("age.gt", 18)` | `WHERE age > 18` |
| `.gte` | `WithFilter("age.gte", 18)` | `WHERE age >= 18` |
| `.lt` | `WithFilter("age.lt", 65)` | `WHERE age < 65` |
| `.lte` | `WithFilter("age.lte", 65)` | `WHERE age <= 65` |
| `.like` | `WithFilter("email.like", "%@gmail.com")` | `WHERE email LIKE '%@gmail.com'` |
| `.null` | `WithFilter("deleted_at.null", true)` | `WHERE deleted_at IS NULL` |
| `.null` | `WithFilter("deleted_at.null", false)` | `WHERE deleted_at IS NOT NULL` |

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
