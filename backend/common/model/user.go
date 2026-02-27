package model

import "gorm.io/gorm"

// User 模型定义 (GORM)
type User struct {
	gorm.Model        // 包含 ID, CreatedAt, UpdatedAt, DeletedAt
	Username   string `gorm:"unique;not null;comment:'用户名'"`
	Password   string `gorm:"not null;comment:'加密后的密码'"`
	Email      string `gorm:"comment:'邮箱'"`
	Role       string `gorm:"default:'user';comment:'角色: user/admin'"`
}

// TableName 指定表名
func (User) TableName() string {
	return "user"
}
