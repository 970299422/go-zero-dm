// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	SQLite struct {
		Path string
	}
	Redis struct {
		Addr      string
		Password  string
		DB        int
		TimeoutMs int
	}
}
