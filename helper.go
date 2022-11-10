package pgxutil

import (
	"context"

	"github.com/allisson/sqlquery"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type FindOptions = sqlquery.FindOptions
type FindAllOptions = sqlquery.FindAllOptions

var (
	postgreSQLFlavor = sqlquery.PostgreSQLFlavor
)

// NewFindOptions returns a FindOptions.
func NewFindOptions() *FindOptions {
	return sqlquery.NewFindOptions(postgreSQLFlavor)
}

// NewFindAllOptions returns a FindAllOptions.
func NewFindAllOptions() *FindAllOptions {
	return sqlquery.NewFindAllOptions(postgreSQLFlavor)
}

// Querier is a abstraction over *pgxpool.Pool/*pgx.Conn/pgx.Tx.
type Querier interface {
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

// Get is a high-level function that calls sqlquery.FindQuery and scany pgxscan.Get function.
func Get(ctx context.Context, db Querier, tableName string, options *FindOptions, dst interface{}) error {
	sqlQuery, args := sqlquery.FindQuery(tableName, options)
	return pgxscan.Get(ctx, db, dst, sqlQuery, args...)
}

// Select is a high-level function that calls sqlquery.FindAllQuery and scany pgxscan.Select function.
func Select(ctx context.Context, db Querier, tableName string, options *FindAllOptions, dst interface{}) error {
	sqlQuery, args := sqlquery.FindAllQuery(tableName, options)
	return pgxscan.Select(ctx, db, dst, sqlQuery, args...)
}

// Insert is a high-level function that calls sqlquery.InsertQuery and pgx Exec.
func Insert(ctx context.Context, db Querier, tag, tableName string, structValue interface{}) error {
	sqlQuery, args := sqlquery.InsertQuery(postgreSQLFlavor, tag, tableName, structValue)
	_, err := db.Exec(ctx, sqlQuery, args...)
	return err
}

// Update is a high-level function that calls sqlquery.pdateQuery and pgx Exec.
func Update(ctx context.Context, db Querier, tag, tableName string, id interface{}, structValue interface{}) error {
	sqlQuery, args := sqlquery.UpdateQuery(postgreSQLFlavor, tag, tableName, id, structValue)
	_, err := db.Exec(ctx, sqlQuery, args...)
	return err
}

// Delete is a high-level function that calls sqlquery.DeleteQuery and pgx Exec.
func Delete(ctx context.Context, db Querier, tableName string, id interface{}) error {
	sqlQuery, args := sqlquery.DeleteQuery(postgreSQLFlavor, tableName, id)
	_, err := db.Exec(ctx, sqlQuery, args...)
	return err
}
