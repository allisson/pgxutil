package pgxutil

import (
	"context"
	"regexp"
	"testing"

	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
)

type player struct {
	ID   int    `db:"id" fieldtag:"insert"`
	Name string `db:"name" fieldtag:"insert,update"`
}

func TestGet(t *testing.T) {
	mock, err := pgxmock.NewConn()
	assert.Nil(t, err)
	defer mock.Close(context.Background())

	rows := mock.NewRows([]string{"id", "name"}).AddRow(1, "Ronaldinho Gaúcho")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM players WHERE id = $1`)).WithArgs("R10").WillReturnRows(rows)

	options := NewFindOptions().WithFilter("id", "R10")
	p := player{}
	err = Get(context.Background(), mock, "players", options, &p)
	assert.Nil(t, err)
	assert.Equal(t, 1, p.ID)
	assert.Equal(t, "Ronaldinho Gaúcho", p.Name)
}

func TestSelect(t *testing.T) {
	mock, err := pgxmock.NewConn()
	assert.Nil(t, err)
	defer mock.Close(context.Background())

	rows := mock.NewRows([]string{"id", "name"}).AddRow(1, "Ronaldinho Gaúcho").AddRow(2, "Ronaldo Fenômeno")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM players LIMIT 10 OFFSET 0 FOR UPDATE SKIP LOCKED`)).WillReturnRows(rows)

	options := NewFindAllOptions().WithLimit(10).WithOffset(0).WithForUpdate("SKIP LOCKED")
	p := []*player{}
	err = Select(context.Background(), mock, "players", options, &p)
	assert.Nil(t, err)
	assert.Len(t, p, 2)
	assert.Equal(t, 1, p[0].ID)
	assert.Equal(t, "Ronaldinho Gaúcho", p[0].Name)
	assert.Equal(t, 2, p[1].ID)
	assert.Equal(t, "Ronaldo Fenômeno", p[1].Name)
}

func TestInsert(t *testing.T) {
	mock, err := pgxmock.NewConn()
	assert.Nil(t, err)
	defer mock.Close(context.Background())

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO players`)).
		WithArgs(1, "Ronaldinho Gaúcho").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	p := player{
		ID:   1,
		Name: "Ronaldinho Gaúcho",
	}
	err = Insert(context.Background(), mock, "", "players", &p)
	assert.Nil(t, err)
}

func TestUpdate(t *testing.T) {
	mock, err := pgxmock.NewConn()
	assert.Nil(t, err)
	defer mock.Close(context.Background())

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE players SET id = $1, name = $2 WHERE id = $3`)).
		WithArgs(1, "Ronaldinho Gaúcho", 1).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	p := player{
		ID:   1,
		Name: "Ronaldinho Gaúcho",
	}
	err = Update(context.Background(), mock, "", "players", p.ID, &p)
	assert.Nil(t, err)
}

func TestDelete(t *testing.T) {
	mock, err := pgxmock.NewConn()
	assert.Nil(t, err)
	defer mock.Close(context.Background())

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM players WHERE id = $1`)).
		WithArgs(1).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	p := player{
		ID:   1,
		Name: "Ronaldinho Gaúcho",
	}
	err = Delete(context.Background(), mock, "players", p.ID)
	assert.Nil(t, err)
}

func TestUpdateWithOptions(t *testing.T) {
	mock, err := pgxmock.NewConn()
	assert.Nil(t, err)
	defer mock.Close(context.Background())

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE players SET name = $1 WHERE id = $2`)).
		WithArgs("Ronaldinho Gaúcho", 1).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	options := NewUpdateOptions().WithAssignment("name", "Ronaldinho Gaúcho").WithFilter("id", 1)
	err = UpdateWithOptions(context.Background(), mock, "players", options)
	assert.Nil(t, err)
}

func TestDeleteWithOptions(t *testing.T) {
	mock, err := pgxmock.NewConn()
	assert.Nil(t, err)
	defer mock.Close(context.Background())

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM players WHERE id = $1`)).
		WithArgs(1).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	options := NewDeleteOptions().WithFilter("id", 1)
	err = DeleteWithOptions(context.Background(), mock, "players", options)
	assert.Nil(t, err)
}
