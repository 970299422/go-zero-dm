// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"time"

	"go-zero-learning/backend/app/identity/api/internal/svc"
	"go-zero-learning/backend/app/identity/api/internal/types"
	"go-zero-learning/backend/common/jwtx"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// todo: add your logic here and delete this line
	// 1. 模拟数据库查询 (Hardcode)
	// TS 类比: if (req.username !== 'admin' || req.password !== '123456') ...
	if req.Username != "admin" || req.Password != "123456" {
		return nil, errors.New("invalid username or password")
	}

	// 2. 定义过期时间 (1小时)
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.JwtAuth.AccessExpire // 1小时
	secretKey := l.svcCtx.Config.JwtAuth.AccessSecret    // // 这里硬编码了密钥！这是一个坏味道(Bad Smell)

	// 3. 生成 Token
	// payload: { "userId": 1 }

	payload := map[string]interface{}{
		"userId": 1,
	}

	token, err := jwtx.GetToken(secretKey, now, accessExpire, payload)
	if err != nil {
		return nil, err
	}

	// 4. 返回响应
	return &types.LoginResp{
		AccessToken:  token,
		AccessExpire: now + accessExpire,
	}, nil
}
