package logic

import (
	"context"

	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCustomerConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCustomerConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCustomerConfigLogic {
	return &DeleteCustomerConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCustomerConfigLogic) DeleteCustomerConfig(req *types.FindCustomerConfigRequest) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
