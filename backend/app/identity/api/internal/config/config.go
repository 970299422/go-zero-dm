// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	DataSource string // 新增
	JwtAuth    struct {
		AccessSecret string
		AccessExpire int64
	}
	IdentityRpc zrpc.RpcClientConf // RPC 客户端配置
}
