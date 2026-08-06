package main

import "testing"

func TestVectorSearch(t *testing.T) {
	db, err := connectSQLiteVec(":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer db.Close()

	var version struct {
		Value string `db:"value"`
	}
	if err := db.NewQuery("SELECT vec_version() AS value").One(&version); err != nil {
		t.Fatalf("query sqlite-vec version: %v", err)
	}
	if version.Value == "" {
		t.Fatal("sqlite-vec returned an empty version")
	}

	statements := []string{
		`CREATE VIRTUAL TABLE vec_items USING vec0(embedding float[3])`,
		`INSERT INTO vec_items(rowid, embedding) VALUES
			(1, '[1, 0, 0]'),
			(2, '[0, 1, 0]'),
			(3, '[0, 0, 1]')`,
	}
	for _, statement := range statements {
		if _, err := db.NewQuery(statement).Execute(); err != nil {
			t.Fatalf("execute %q: %v", statement, err)
		}
	}

	var nearest struct {
		RowID int `db:"rowid"`
	}
	err = db.NewQuery(`
		SELECT rowid
		FROM vec_items
		WHERE embedding MATCH '[0.9, 0.1, 0]'
		ORDER BY distance
		LIMIT 1
	`).One(&nearest)
	if err != nil {
		t.Fatalf("run nearest-neighbor query: %v", err)
	}
	if nearest.RowID != 1 {
		t.Fatalf("nearest row = %d, want 1", nearest.RowID)
	}
}
