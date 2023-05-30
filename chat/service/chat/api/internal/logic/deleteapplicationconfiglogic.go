package logic

import (
	"chat/service/chat/api/internal/repository"
	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteApplicationConfigLogic struct {
	logx.Logger
	ctx                         context.Context
	svcCtx                      *svc.ServiceContext
	applicationConfigRepository *repository.ApplicationConfigRepository
}

func NewDeleteApplicationConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteApplicationConfigLogic {
	return &DeleteApplicationConfigLogic{
		Logger:                      logx.WithContext(ctx),
		ctx:                         ctx,
		svcCtx:                      svcCtx,
		applicationConfigRepository: repository.NewApplicationConfigRepository(ctx, svcCtx),
	}
}

// TODO 先不支持
func (l *DeleteApplicationConfigLogic) DeleteApplicationConfig(req *types.FindApplicationConfigRequest) (resp *types.Response, err error) {

	return &types.Response{
		Message: "ok",
	}, nil
	return
}
