package dbtypegen_test

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	dbtype "github.com/johejo/dbtypegen/testdata"
)

func Example() {
	ctx := context.Background()
	db, err := sql.Open("mysql", "root:pass@tcp(localhost:3306)/dbtypegen?charset=utf8mb4&parseTime=true&loc=UTC&multiStatements=true")
	if err != nil {
		panic(err)
	}
	defer func() {
		if _, err := db.ExecContext(ctx, "DROP TABLE `user`"); err != nil {
			panic(err)
		}
		if _, err := db.ExecContext(ctx, "DROP TABLE `group`"); err != nil {
			panic(err)
		}
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()

	schema, err := ioutil.ReadFile("testdata/schema.sql")
	if err != nil {
		panic(err)
	}

	if _, err := db.ExecContext(ctx, string(schema)); err != nil {
		panic(err)
	}

	var u dbtype.User

	var b strings.Builder
	b.WriteString("INSERT INTO ")
	b.WriteString(u.TableName())
	b.WriteString(" (")
	b.WriteString(u.Columns())
	b.WriteString(") ")
	b.WriteString("VALUES (?,?,?,?)")
	q := b.String()
	now := time.Now()
	args := []interface{}{1, now, true, "Gopher"}

	if _, err := db.ExecContext(ctx, q, args...); err != nil {
		panic(err)
	}

	if err := db.QueryRowContext(ctx, u.SelectAll()+" WHERE id=?", 1).Scan(u.Scans()...); err != nil {
		panic(err)
	}

	fmt.Println(u)
}
