// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"context"
	"fmt"
	"time"

	"go-zero-learning/backend/app/demo/api/internal/config"
	"go-zero-learning/backend/app/demo/api/internal/seed"
	"go-zero-learning/backend/app/demo/api/internal/store"

	"github.com/redis/go-redis/v9"
)

type ServiceContext struct {
	Config config.Config
	Store  *store.SQLiteStore
	Redis  *redis.Client
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

	timeout := time.Duration(c.Redis.TimeoutMs) * time.Millisecond

	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		Password:     c.Redis.Password,
		DB:           c.Redis.DB,
		DialTimeout:  timeout,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		MaxRetries:   0,
	})

	return &ServiceContext{
		Config: c,
		Store:  sqliteStore,
		Redis:  rdb,
	}
}

func (s *ServiceContext) Close() error {
	if s == nil || s.Store == nil {
		return nil
	}

	return s.Store.Close()
}
