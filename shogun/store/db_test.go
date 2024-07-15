package store_test

import (
	"context"
	gosql "database/sql"
	"fmt"
	"testing"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"

	"github.com/nwillems/shogun-graph/shogun/store"
	"github.com/nwillems/shogun-graph/shogun/store/sql"
)

//go:embed sql/migration/db_up.sql
var ddl string

func TestStorage(t *testing.T) {
	// Init
	ctx := context.Background()
	db, err := gosql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	s := sql.NewStore(db)
	err = s.RunMigration(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Test
	err = s.Store(ctx,
		[]store.CreateNode{
			{Name: "A"}, {Name: "B"}, {Name: "C"}, {Name: "D"}, {Name: "E"},
		},
		[]store.CreateEdge{
			{Source: "A", Target: "B", Type: "lnk"},
			{Source: "A", Target: "D", Type: "lnk"},
			{Source: "A", Target: "C", Type: "lnk"},
			{Source: "B", Target: "A", Type: "lnk"},
			{Source: "C", Target: "D", Type: "lnk"},
			{Source: "D", Target: "A", Type: "lnk"},
			{Source: "D", Target: "C", Type: "lnk"},
			{Source: "D", Target: "E", Type: "lnk"},
			{Source: "E", Target: "C", Type: "lnk"},
		})
	if err != nil {
		t.Fatal(err)
	}

	/*
		res, err := s.Query(store.Query{
			Nodes: []string{"A", "E"},
			})
			if err != nil {
				t.Fatal(err)
				}*/

	res, _ := db.Query("SELECT node_name FROM nodes")
	for res.Next() {
		var name string
		_ = res.Scan(&name)
		fmt.Printf("%+v\n", name)
	}

	res, _ = db.Query("SELECT source,target FROM edges")
	for res.Next() {
		var src, dst string
		_ = res.Scan(&src, &dst)
		fmt.Printf("%+v -> %+v\n", src, dst)
	}
	fmt.Printf("%+v\n", res)

	rr, err := s.QueryNode(ctx, 1)
	fmt.Printf("%+v\n", rr)

	t.Fatal("THIS SHOULD NOT FAIL")
}
