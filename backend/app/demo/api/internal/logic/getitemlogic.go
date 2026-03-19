// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-zero-learning/backend/app/demo/api/internal/svc"
	"go-zero-learning/backend/app/demo/api/internal/types"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetItemLogic {
	return &GetItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetItemLogic) GetItem(req *types.GetItemReq) (resp *types.ItemResp, err error) {
	key := fmt.Sprintf("demo:item:%d", req.Id)

	cacheCtx, cancel := context.WithTimeout(l.ctx, time.Duration(l.svcCtx.Config.Redis.TimeoutMs)*time.Millisecond)
	defer cancel()

	cached, cacheErr := l.svcCtx.Redis.Get(cacheCtx, key).Result()
	if cacheErr == nil {
		var cachedResp types.ItemResp
		if err := json.Unmarshal([]byte(cached), &cachedResp); err == nil {
			return &cachedResp, nil
		}
	}

	if cacheErr != nil && cacheErr != redis.Nil {
		logx.Errorf("redis get failed: %v", cacheErr)
	}

	item, err := l.svcCtx.Store.GetItem(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}

	resp = &types.ItemResp{
		Id:        item.ID,
		Name:      item.Name,
		UpdatedAt: item.UpdatedAt,
	}

	raw, err := json.Marshal(resp)
	if err == nil {
		setCtx, setCancel := context.WithTimeout(l.ctx, time.Duration(l.svcCtx.Config.Redis.TimeoutMs)*time.Millisecond)
		defer setCancel()
		if setErr := l.svcCtx.Redis.Set(setCtx, key, raw, 60*time.Second).Err(); setErr != nil {
			logx.Errorf("redis set failed: %v", setErr)
		}
	}

	return resp, nil
}
