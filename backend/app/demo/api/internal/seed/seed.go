package seed

import (
	"context"

	"go-zero-learning/backend/app/demo/api/internal/store"
)

var defaultItems = []store.Item{
	{
		ID:        1,
		Name:      "Hello, SQLite",
		UpdatedAt: "2026-03-18T00:00:00Z",
	},
	{
		ID:        2,
		Name:      "Start here",
		UpdatedAt: "2026-03-18T00:00:00Z",
	},
}

func Seed(ctx context.Context, store *store.SQLiteStore) error {
	for _, item := range defaultItems {
		if err := store.InsertItemIfMissing(ctx, item); err != nil {
			return err
		}
	}

	return nil
}
