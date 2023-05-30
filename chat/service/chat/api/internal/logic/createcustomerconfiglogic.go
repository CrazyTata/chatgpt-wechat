package logic

import (
	"context"

	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCustomerConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCustomerConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCustomerConfigLogic {
	return &CreateCustomerConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCustomerConfigLogic) CreateCustomerConfig(req *types.CustomerConfig) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
