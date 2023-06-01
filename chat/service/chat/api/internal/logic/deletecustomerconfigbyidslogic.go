package logic

import (
	"chat/common/util"
	"chat/service/chat/api/internal/repository"
	"context"

	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCustomerConfigByIdsLogic struct {
	logx.Logger
	ctx                      context.Context
	svcCtx                   *svc.ServiceContext
	customerConfigRepository *repository.CustomerConfigRepository
}

func NewDeleteCustomerConfigByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCustomerConfigByIdsLogic {
	return &DeleteCustomerConfigByIdsLogic{
		Logger:                   logx.WithContext(ctx),
		ctx:                      ctx,
		svcCtx:                   svcCtx,
		customerConfigRepository: repository.NewCustomerConfigRepository(ctx, svcCtx),
	}
}

func (l *DeleteCustomerConfigByIdsLogic) DeleteCustomerConfigByIds(req *types.IdsRequest) (resp *types.Response, err error) {
	if len(req.Ids) > 0 {
		for _, v := range req.Ids {
			err = l.customerConfigRepository.Delete(v)
			if err != nil {
				util.Info("DeleteCustomerConfigByIds err:" + err.Error())
			}
		}
	}

	return &types.Response{
		Message: "ok",
	}, nil
}
