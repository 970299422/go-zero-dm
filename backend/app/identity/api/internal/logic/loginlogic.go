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
	"go-zero-learning/backend/common/model"

	"golang.org/x/crypto/bcrypt"

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
	// 1. 根据用户名查询用户
	var user model.User
	// GORM: Where条件查询，First取第一条
	result := l.svcCtx.DB.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		return nil, errors.New("用户不存在或登录失败")
	}

	// 2. 校验密码 (Hash对比)
	// CompareHashAndPassword(数据库里的Hash值, 用户输入的明文密码)
	// 如果不匹配，会返回 error
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 3. 密码正确，生成 JWT Token
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.JwtAuth.AccessExpire //这是有效期时长(秒)
	secretKey := l.svcCtx.Config.JwtAuth.AccessSecret

	// Payload: 放入最重要的 userId，以便后续接口知道是谁在操作
	payload := map[string]interface{}{
		"userId": user.ID,
	}

	token, err := jwtx.GetToken(secretKey, now, accessExpire, payload)
	if err != nil {
		return nil, err
	}

	return &types.LoginResp{
		AccessToken:  token,
		AccessExpire: now + accessExpire, // 返回绝对时间戳
	}, nil
}
