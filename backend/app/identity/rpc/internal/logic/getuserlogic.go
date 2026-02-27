package logic

import (
	"context"
	"errors"

	"go-zero-learning/backend/app/identity/rpc/internal/svc"
	pb "go-zero-learning/backend/app/identity/rpc/pb"
	"go-zero-learning/backend/common/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 定义一个 RPC 方法：传入 GetUserReq，返回 GetUserResp
func (l *GetUserLogic) GetUser(in *pb.GetUserReq) (*pb.GetUserResp, error) {
	// 1. 查询数据库
	var user model.User
	result := l.svcCtx.DB.First(&user, in.Id)
	if result.Error != nil {
		return nil, errors.New("RPC错误: 用户不存在")
	}

	// 2. 转换为 pb 对象返回
	return &pb.GetUserResp{
		Id:       int64(user.ID),
		Username: user.Username,
	}, nil
}
