// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"time"

	"go-zero-learning/backend/app/identity/api/internal/svc"
	"go-zero-learning/backend/app/identity/api/internal/types"
	"go-zero-learning/backend/app/identity/rpc/identity"
	"go-zero-learning/backend/common/jwtx"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	// 调用 RPC 的 CreateUser 方法（gRPC）
	rpcResp, err := l.svcCtx.IdentityRpc.CreateUser(l.ctx, &identity.CreateUserReq{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, errors.New("注册失败: " + err.Error())
	}

	// 4. 生成 JWT Token (注册成功后直接登录)
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.JwtAuth.AccessExpire
	secretKey := l.svcCtx.Config.JwtAuth.AccessSecret

	// Payload: 放入最重要的 userId，以便后续接口知道是谁在操作
	payload := map[string]interface{}{
		"userId": rpcResp.Id,
	}

	token, err := jwtx.GetToken(secretKey, now, accessExpire, payload)
	if err != nil {
		return nil, err
	}

	return &types.RegisterResp{
		AccessToken:  token,
		AccessExpire: now + accessExpire,
	}, nil
}
