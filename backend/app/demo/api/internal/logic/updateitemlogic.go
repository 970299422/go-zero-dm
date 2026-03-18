// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"go-zero-learning/backend/app/demo/api/internal/svc"
	"go-zero-learning/backend/app/demo/api/internal/store"
	"go-zero-learning/backend/app/demo/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateItemLogic {
	return &UpdateItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateItemLogic) UpdateItem(req *types.UpdateItemReq) (resp *types.ItemResp, err error) {
	item := store.Item{
		ID:        req.Id,
		Name:      req.Name,
		UpdatedAt: "2026-03-18T00:00:00Z",
	}
	if err := l.svcCtx.Store.SaveItem(l.ctx, item); err != nil {
		return nil, err
	}

	return &types.ItemResp{
		Id:        item.ID,
		Name:      item.Name,
		UpdatedAt: item.UpdatedAt,
	}, nil
}
