// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"flag"
	"fmt"

	"go-zero-learning/backend/app/identity/api/internal/config"
	"go-zero-learning/backend/app/identity/api/internal/handler"
	"go-zero-learning/backend/app/identity/api/internal/middleware"
	"go-zero-learning/backend/app/identity/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/identity-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 注册全局轻量 Token 提取中间件（不会替代 go-zero 的鉴权）
	server.Use(middleware.JWTMiddleware(c.JwtAuth.AccessSecret))

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
