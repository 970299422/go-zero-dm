package logic

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"go-zero-learning/backend/app/identity/rpc/internal/svc"
	"go-zero-learning/backend/app/identity/rpc/pb"
	"go-zero-learning/backend/common/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateUserLogic) CreateUser(in *pb.CreateUserReq) (*pb.CreateUserResp, error) {
	// 查询用户名是否存在
	var existing model.User
	if err := l.svcCtx.DB.Where("username = ?", in.Username).First(&existing).Error; err == nil {
		return nil, ErrUsernameExists
	}

	// 加密密码
	hashed, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		l.Logger.Error(err)
		return nil, err
	}

	user := model.User{
		Username: in.Username,
		Password: string(hashed),
	}

	if err := l.svcCtx.DB.Create(&user).Error; err != nil {
		l.Logger.Error(err)
		return nil, err
	}

	return &pb.CreateUserResp{
		Id:       int64(user.ID),
		Username: user.Username,
	}, nil
}

// ErrUsernameExists is a sentinel error used to indicate duplicate username.
var ErrUsernameExists = errors.New("用户名已存在")
