package logic

import (
	"context"

	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCustomerConfigByIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCustomerConfigByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCustomerConfigByIdsLogic {
	return &DeleteCustomerConfigByIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCustomerConfigByIdsLogic) DeleteCustomerConfigByIds(req *types.IdsRequest) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
