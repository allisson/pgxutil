// Package pgxutil provides high-level helper functions for PostgreSQL database operations
// using pgx driver. It wraps sqlquery for SQL generation and scany for row scanning,
// offering convenient methods for common CRUD operations with type safety.
package pgxutil

import (
	"context"

	"github.com/allisson/sqlquery"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type (
	// FindOptions configures query options for finding a single record.
	// It supports filtering, column selection, and ordering for SELECT queries.
	FindOptions = sqlquery.FindOptions

	// FindAllOptions configures query options for finding multiple records.
	// It supports filtering, column selection, ordering, and pagination for SELECT queries.
	FindAllOptions = sqlquery.FindAllOptions

	// UpdateOptions configures query options for updating records.
	// It supports setting values, filtering conditions, and returning clauses.
	UpdateOptions = sqlquery.UpdateOptions

	// DeleteOptions configures query options for deleting records.
	// It supports filtering conditions and returning clauses.
	DeleteOptions = sqlquery.DeleteOptions
)

var (
	postgreSQLFlavor = sqlquery.PostgreSQLFlavor
)

// NewFindOptions returns a FindOptions configured for PostgreSQL.
func NewFindOptions() *FindOptions {
	return sqlquery.NewFindOptions(postgreSQLFlavor)
}

// NewFindAllOptions returns a FindAllOptions configured for PostgreSQL.
func NewFindAllOptions() *FindAllOptions {
	return sqlquery.NewFindAllOptions(postgreSQLFlavor)
}

// NewUpdateOptions returns an UpdateOptions configured for PostgreSQL.
func NewUpdateOptions() *UpdateOptions {
	return sqlquery.NewUpdateOptions(postgreSQLFlavor)
}

// NewDeleteOptions returns DeleteOptions configured for PostgreSQL.
func NewDeleteOptions() *DeleteOptions {
	return sqlquery.NewDeleteOptions(postgreSQLFlavor)
}

// Querier is an abstraction over *pgxpool.Pool, *pgx.Conn, and pgx.Tx.
// It allows functions to work with any of these pgx types for query execution.
type Querier interface {
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

// Get retrieves a single record from the database and scans it into dst.
// It generates a SELECT query based on the provided options and returns an error if no record is found.
func Get(ctx context.Context, db Querier, tableName string, options *FindOptions, dst interface{}) error {
	sqlQuery, args := sqlquery.FindQuery(tableName, options)
	return pgxscan.Get(ctx, db, dst, sqlQuery, args...)
}

// Select retrieves multiple records from the database and scans them into dst.
// It generates a SELECT query based on the provided options. The dst parameter should be a pointer to a slice.
func Select(ctx context.Context, db Querier, tableName string, options *FindAllOptions, dst interface{}) error {
	sqlQuery, args := sqlquery.FindAllQuery(tableName, options)
	return pgxscan.Select(ctx, db, dst, sqlQuery, args...)
}

// Insert inserts a new record into the database.
// It generates an INSERT query from the struct fields based on the provided tag (e.g., "db").
func Insert(ctx context.Context, db Querier, tag, tableName string, structValue interface{}) error {
	sqlQuery, args := sqlquery.InsertQuery(postgreSQLFlavor, tag, tableName, structValue)
	_, err := db.Exec(ctx, sqlQuery, args...)
	return err
}

// Update updates a record in the database by its ID.
// It generates an UPDATE query from the struct fields based on the provided tag (e.g., "db").
func Update(ctx context.Context, db Querier, tag, tableName string, id interface{}, structValue interface{}) error {
	sqlQuery, args := sqlquery.UpdateQuery(postgreSQLFlavor, tag, tableName, id, structValue)
	_, err := db.Exec(ctx, sqlQuery, args...)
	return err
}

// Delete deletes a record from the database by its ID.
// It generates a DELETE query with a WHERE clause matching the provided ID.
func Delete(ctx context.Context, db Querier, tableName string, id interface{}) error {
	sqlQuery, args := sqlquery.DeleteQuery(postgreSQLFlavor, tableName, id)
	_, err := db.Exec(ctx, sqlQuery, args...)
	return err
}

// UpdateWithOptions updates records in the database using custom options.
// It allows fine-grained control over which columns to update and which records to match.
func UpdateWithOptions(ctx context.Context, db Querier, tableName string, options *UpdateOptions) error {
	sqlQuery, args := sqlquery.UpdateWithOptionsQuery(tableName, options)
	_, err := db.Exec(ctx, sqlQuery, args...)
	return err
}

// DeleteWithOptions deletes records from the database using custom options.
// It allows fine-grained control over which records to delete with complex WHERE conditions.
func DeleteWithOptions(ctx context.Context, db Querier, tableName string, options *DeleteOptions) error {
	sqlQuery, args := sqlquery.DeleteWithOptionsQuery(tableName, options)
	_, err := db.Exec(ctx, sqlQuery, args...)
	return err
}
