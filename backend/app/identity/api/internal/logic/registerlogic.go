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
	// 1. 检查用户名是否已存在
	var existingUser model.User
	// GORM 查询: Where("条件", 参数).First(&对象)
	// RowsAffected > 0 表示找到了记录
	if result := l.svcCtx.DB.Where("username = ?", req.Username).First(&existingUser); result.RowsAffected > 0 {
		return nil, errors.New("用户名已存在")
	}

	// 2. 密码加密 (Hash)
	// 使用 bcrypt 算法。它会自动加盐(salt)，非常安全。
	// GenerateFromPassword 接受此时的明文密码和代价(Cost)。DefaultCost 即可。
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("系统错误: 密码加密失败")
	}

	// 3. 创建新用户
	newUser := model.User{
		Username: req.Username,
		Password: string(hashedPassword), // 存入数据库的是加密后的乱码
	}

	// GORM 插入: Create(&对象)
	if result := l.svcCtx.DB.Create(&newUser); result.Error != nil {
		return nil, errors.New("注册失败: " + result.Error.Error())
	}

	// 4. 生成 JWT Token (注册成功后直接登录)
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.JwtAuth.AccessExpire
	secretKey := l.svcCtx.Config.JwtAuth.AccessSecret

	// Payload: 放入最重要的 userId，以便后续接口知道是谁在操作
	payload := map[string]interface{}{
		"userId": newUser.ID,
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
