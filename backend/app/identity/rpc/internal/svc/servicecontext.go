package svc

import (
	"go-zero-learning/backend/app/identity/rpc/internal/config"
	// 这里为了简化，我们直接引用 API 层的 model
	// "go-zero-learning/backend/common/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB // 全局数据库连接对象
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 连接 SQLite 数据库 (和 API 一样)
	db, err := gorm.Open(sqlite.Open(c.DataSource), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	return &ServiceContext{
		Config: c,
		DB:     db,
	}
}
