package logic

import (
	"context"

	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteApplicationConfigByIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteApplicationConfigByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteApplicationConfigByIdsLogic {
	return &DeleteApplicationConfigByIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteApplicationConfigByIdsLogic) DeleteApplicationConfigByIds(req *types.IdsRequest) (resp *types.Response, err error) {
	// TODO 先不支持
	return
}
