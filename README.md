# pgxutil
[![Build Status](https://github.com/allisson/pgxutil/workflows/Release/badge.svg)](https://github.com/allisson/pgxutil/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/allisson/pgxutil)](https://goreportcard.com/report/github.com/allisson/pgxutil)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/allisson/pgxutil)

A collection of helpers to deal with pgx toolkit.

Example:

```golang
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
}

func main() {
	ctx := context.Background()
	// Run a database with docker: docker run --name test --restart unless-stopped -e POSTGRES_USER=user -e POSTGRES_PASSWORD=password -e POSTGRES_DB=pgxutil -p 5432:5432 -d postgres:14-alpine
	// Connect to database
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
			age INTEGER NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Insert players
	r9 := Player{
		Name: "Ronaldo Fenômeno",
		Age:  44,
	}
	r10 := Player{
		Name: "Ronaldinho Gaúcho",
		Age:  41,
	}
	tag := "insert" // will use fields with fieldtag:"insert"
	if err := pgxutil.Insert(ctx, conn, tag, "players", &r9); err != nil {
		log.Fatal(err)
	}
	if err := pgxutil.Insert(ctx, conn, tag, "players", &r10); err != nil {
		log.Fatal(err)
	}

	// Get player
	findOptions := pgxutil.NewFindOptions().WithFilter("name", r10.Name)
	if err := pgxutil.Get(ctx, conn, "players", findOptions, &r10); err != nil {
		log.Fatal(err)
	}
	findOptions = pgxutil.NewFindOptions().WithFilter("name", r9.Name)
	if err := pgxutil.Get(ctx, conn, "players", findOptions, &r9); err != nil {
		log.Fatal(err)
	}

	// Select players
	players := []*Player{}
	findAllOptions := pgxutil.NewFindAllOptions().WithLimit(10).WithOffset(0).WithOrderBy("name asc")
	if err := pgxutil.Select(ctx, conn, "players", findAllOptions, &players); err != nil {
		log.Fatal(err)
	}
	for _, p := range players {
		fmt.Printf("%#v\n", p)
	}

	// Update player
	tag = "update" // will use fields with fieldtag:"update"
	r10.Name = "Ronaldinho Bruxo"
	if err := pgxutil.Update(ctx, conn, tag, "players", r10.ID, &r10); err != nil {
		log.Fatal(err)
	}

	// Delete player
	if err := pgxutil.Delete(ctx, conn, "players", r9.ID); err != nil {
		log.Fatal(err)
	}
}
```

Options for FindOptions and FindAllOptions:

```golang
package main

import (
	"github.com/allisson/pgxutil/v2"
)

func main() {
	findOptions := pgxutil.NewFindOptions().
		WithFields([]string{"id", "name"}). // Return only id and name fields
		WithFilter("id", 1).                // WHERE id = 1
		WithFilter("id", nil).              // WHERE id IS NULL
		WithFilter("id.in", "1,2,3").       // WHERE id IN (1, 2, 3)
		WithFilter("id.notin", "1,2,3").    // WHERE id NOT IN ($1, $2, $3)
		WithFilter("id.not", 1).            // WHERE id <> 1
		WithFilter("id.gt", 1).             // WHERE id > 1
		WithFilter("id.gte", 1).            // WHERE id >= 1
		WithFilter("id.lt", 1).             // WHERE id < 1
		WithFilter("id.lte", 1).            // WHERE id <= 1
		WithFilter("id.like", 1).           // WHERE id LIKE 1
		WithFilter("id.null", true).        // WHERE id.null IS NULL
		WithFilter("id.null", false)        // WHERE id.null IS NOT NULL

	findAllOptions := pgxutil.NewFindAllOptions().
		WithFields([]string{"id", "name"}). // Return only id and name fields
		WithFilter("id", 1).                // WHERE id = 1
		WithFilter("id", nil).              // WHERE id IS NULL
		WithFilter("id.in", "1,2,3").       // WHERE id IN (1, 2, 3)
		WithFilter("id.notin", "1,2,3").    // WHERE id NOT IN ($1, $2, $3)
		WithFilter("id.not", 1).            // WHERE id <> 1
		WithFilter("id.gt", 1).             // WHERE id > 1
		WithFilter("id.gte", 1).            // WHERE id >= 1
		WithFilter("id.lt", 1).             // WHERE id < 1
		WithFilter("id.lte", 1).            // WHERE id <= 1
		WithFilter("id.like", 1).           // WHERE id LIKE 1
		WithFilter("id.null", true).        // WHERE id.null IS NULL
		WithFilter("id.null", false).       // WHERE id.null IS NOT NULL
		WithLimit(10).                      // LIMIT 10
		WithOffset(0).                      // OFFSET 0
		WithOrderBy("name asc")             // ORDER BY name asc
}
```
