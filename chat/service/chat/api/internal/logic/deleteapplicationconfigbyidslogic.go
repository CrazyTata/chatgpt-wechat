package logic

import (
	"chat/common/util"
	"chat/service/chat/api/internal/repository"
	"context"

	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteApplicationConfigByIdsLogic struct {
	logx.Logger
	ctx                         context.Context
	svcCtx                      *svc.ServiceContext
	applicationConfigRepository *repository.ApplicationConfigRepository
}

func NewDeleteApplicationConfigByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteApplicationConfigByIdsLogic {
	return &DeleteApplicationConfigByIdsLogic{
		Logger:                      logx.WithContext(ctx),
		ctx:                         ctx,
		svcCtx:                      svcCtx,
		applicationConfigRepository: repository.NewApplicationConfigRepository(ctx, svcCtx),
	}
}

func (l *DeleteApplicationConfigByIdsLogic) DeleteApplicationConfigByIds(req *types.IdsRequest) (resp *types.Response, err error) {
	if len(req.Ids) > 0 {
		for _, v := range req.Ids {
			err = l.applicationConfigRepository.Delete(v)
			if err != nil {
				util.Info("applicationConfigRepository err:" + err.Error())
			}
		}
	}

	return &types.Response{
		Message: "ok",
	}, nil
}
