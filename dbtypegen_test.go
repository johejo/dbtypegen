package dbtypegen

import (
	"context"
	"database/sql"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/go-cmp/cmp"

	dbtype "github.com/johejo/dbtypegen/testdata"
)

func TestGenerate(t *testing.T) {
	opts := []Option{}

	ctx := context.Background()
	f, err := os.Open("testdata/schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	generated, err := Generate(ctx, f, opts...)
	if err != nil {
		t.Fatal(err)
	}

	golden, err := ioutil.ReadFile("testdata/golden.go")
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(string(generated), string(golden)); diff != "" {
		t.Error(diff)
	}
	t.Log(string(generated))
}

func TestGenerated(t *testing.T) {
	ctx := context.Background()
	db, err := sql.Open("mysql", "root:pass@tcp(localhost:3306)/dbtypegen?charset=utf8mb4&parseTime=true&loc=UTC&multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := db.ExecContext(ctx, "DROP TABLE `user`"); err != nil {
			t.Fatal(err)
		}
		if _, err := db.ExecContext(ctx, "DROP TABLE `group`"); err != nil {
			t.Fatal(err)
		}
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	})

	schema, err := ioutil.ReadFile("testdata/schema.sql")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := db.ExecContext(ctx, string(schema)); err != nil {
		t.Fatal(err)
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
	now := time.Now().In(time.UTC).Round(1 * time.Millisecond)
	args := []interface{}{1, now, true, "Gopher"}

	if _, err := db.ExecContext(ctx, q, args...); err != nil {
		t.Fatal(err)
	}

	if err := db.QueryRowContext(ctx, u.SelectAll()+" WHERE id=?", 1).Scan(u.Scans()...); err != nil {
		t.Fatal(err)
	}

	want := dbtype.User{Id: 1, CreatedAt: now, Active: true, Name: "Gopher"}
	if diff := cmp.Diff(u, want); diff != "" {
		t.Fatal(diff)
	}
}
