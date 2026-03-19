package svc

import (
	"context"
	"path/filepath"
	"testing"

	"go-zero-learning/backend/app/demo/api/internal/config"
)

func TestNewServiceContextSeedsSQLite(t *testing.T) {
	t.Parallel()

	c := config.Config{}
	c.SQLite.Path = filepath.Join(t.TempDir(), "demo.db")

	ctx := NewServiceContext(c)
	defer func() {
		_ = ctx.Close()
	}()

	first, err := ctx.Store.GetItem(context.Background(), 1)
	if err != nil {
		t.Fatalf("get seeded item 1: %v", err)
	}
	if first.ID != 1 || first.Name != "Hello, SQLite" || first.UpdatedAt != "2026-03-18T00:00:00Z" {
		t.Fatalf("unexpected seeded item 1: %+v", first)
	}

	second, err := ctx.Store.GetItem(context.Background(), 2)
	if err != nil {
		t.Fatalf("get seeded item 2: %v", err)
	}
	if second.ID != 2 || second.Name != "Start here" || second.UpdatedAt != "2026-03-18T00:00:00Z" {
		t.Fatalf("unexpected seeded item 2: %+v", second)
	}
}
