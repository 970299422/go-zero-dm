// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"go-zero-learning/backend/app/identity/api/internal/config"
	"go-zero-learning/backend/app/identity/api/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB // 全局数据库连接对象
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 连接 SQLite 数据库
	db, err := gorm.Open(sqlite.Open(c.DataSource), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// 自动建表 (Auto Migrate)
	// 每次启动服务时，GORM 会检查表结构变化并自动同步
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	return &ServiceContext{
		Config: c,
		DB:     db,
	}
}
