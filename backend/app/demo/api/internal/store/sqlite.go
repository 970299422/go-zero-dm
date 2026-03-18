package store

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Item struct {
	ID        int64
	Name      string
	UpdatedAt string
}

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(path string) (*SQLiteStore, error) {
	if path == "" {
		return nil, fmt.Errorf("sqlite path is required")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create sqlite directory: %w", err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	db.SetMaxOpenConns(1)

	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Close() error {
	if s == nil || s.db == nil {
		return nil
	}

	return s.db.Close()
}

func (s *SQLiteStore) InitSchema(ctx context.Context) error {
	const schema = `
CREATE TABLE IF NOT EXISTS items (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	updated_at TEXT NOT NULL
);
`

	if _, err := s.db.ExecContext(ctx, schema); err != nil {
		return fmt.Errorf("create items table: %w", err)
	}

	return nil
}

func (s *SQLiteStore) GetItem(ctx context.Context, id int64) (Item, error) {
	var item Item

	row := s.db.QueryRowContext(ctx, `
SELECT id, name, updated_at
FROM items
WHERE id = ?;
`, id)
	if err := row.Scan(&item.ID, &item.Name, &item.UpdatedAt); err != nil {
		return Item{}, fmt.Errorf("query item %d: %w", id, err)
	}

	return item, nil
}

func (s *SQLiteStore) InsertItemIfMissing(ctx context.Context, item Item) error {
	_, err := s.db.ExecContext(ctx, `
INSERT OR IGNORE INTO items (id, name, updated_at)
VALUES (?, ?, ?);
`, item.ID, item.Name, item.UpdatedAt)
	if err != nil {
		return fmt.Errorf("seed item %d: %w", item.ID, err)
	}

	return nil
}

func (s *SQLiteStore) SaveItem(ctx context.Context, item Item) error {
	_, err := s.db.ExecContext(ctx, `
INSERT INTO items (id, name, updated_at)
VALUES (?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
	name = excluded.name,
	updated_at = excluded.updated_at;
`, item.ID, item.Name, item.UpdatedAt)
	if err != nil {
		return fmt.Errorf("save item %d: %w", item.ID, err)
	}

	return nil
}
