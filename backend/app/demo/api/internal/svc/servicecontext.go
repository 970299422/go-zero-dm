// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"context"
	"fmt"

	"go-zero-learning/backend/app/demo/api/internal/config"
	"go-zero-learning/backend/app/demo/api/internal/seed"
	"go-zero-learning/backend/app/demo/api/internal/store"
)

type ServiceContext struct {
	Config config.Config
	Store  *store.SQLiteStore
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqliteStore, err := store.NewSQLiteStore(c.SQLite.Path)
	if err != nil {
		panic(fmt.Sprintf("failed to open sqlite database: %v", err))
	}

	if err := sqliteStore.InitSchema(context.Background()); err != nil {
		panic(fmt.Sprintf("failed to initialize sqlite schema: %v", err))
	}

	if err := seed.Seed(context.Background(), sqliteStore); err != nil {
		panic(fmt.Sprintf("failed to seed sqlite database: %v", err))
	}

	return &ServiceContext{
		Config: c,
		Store:  sqliteStore,
	}
}
