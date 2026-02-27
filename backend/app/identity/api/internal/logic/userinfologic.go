// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"encoding/json"
	"errors"

	"go-zero-learning/backend/app/identity/api/internal/svc"
	"go-zero-learning/backend/app/identity/api/internal/types"
	"go-zero-learning/backend/app/identity/rpc/identity"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo() (resp *types.UserInfoResp, err error) {
	// 1. 从 Context 中获取 userId
	// 注意：go-zero 的 JWT 中间件解析 Token 后，把 payload 里的值放进了 ctx
	// 还没完！取出来的 value 类型很可能是 json.Number (本质是 string)
	userIdVal := l.ctx.Value("userId")
	if userIdVal == nil {
		return nil, errors.New("无效的Token: 找不到userId")
	}

	// 2. 类型安全转换: json.Number -> int64
	var userId int64
	// jwt 解析出来的数字通常是 json.Number 类型
	if v, ok := userIdVal.(json.Number); ok {
		var err error
		userId, err = v.Int64()
		if err != nil {
			return nil, errors.New("无效的Token: userId格式错误")
		}
	} else if v, ok := userIdVal.(float64); ok { // 防止某些版本解析为 float64
		userId = int64(v)
	} else if v, ok := userIdVal.(int64); ok {
		userId = v
	} else {
		// 打印一下实际类型，方便调试
		l.Errorf("userId type mismatch: %T", userIdVal)
		return nil, errors.New("无效的Token: userId类型异常")
	}

	// 3. 调用 RPC 服务获取用户信息 (不再直接查数据库)
	// 引入 pb 包: "go-zero-learning/backend/app/identity/rpc/pb"
	// 注意：这里把 l.ctx 传给了 RPC，这样如果 API 超时，RPC 也会自动取消
	rpcResp, err := l.svcCtx.IdentityRpc.GetUser(l.ctx, &identity.GetUserReq{
		Id: userId,
	})
	if err != nil {
		return nil, errors.New("获取用户信息失败: " + err.Error())
	}

	// 4. 返回给前端
	return &types.UserInfoResp{
		Id:       rpcResp.Id,
		Username: rpcResp.Username,
	}, nil
}
