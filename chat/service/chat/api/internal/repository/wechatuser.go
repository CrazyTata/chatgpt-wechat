package repository

import (
	"chat/service/chat/api/internal/svc"
	"chat/service/chat/model"
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/logx"
)

type WechatUserRepository struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWechatUserRepository(ctx context.Context, svcCtx *svc.ServiceContext) *WechatUserRepository {
	return &WechatUserRepository{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WechatUserRepository) GetByName(kfName string) (customerPo *model.WechatUser, err error) {

	customerPo, err = l.svcCtx.WechatUserModel.FindOneByQuery(context.Background(),
		l.svcCtx.WechatUserModel.RowBuilder().Where(squirrel.Eq{"nickname": kfName}),
	)
	return
}

func (l *WechatUserRepository) GetByUsers(users []string) (CustomerConfigPo []*model.WechatUser, err error) {
	return l.svcCtx.WechatUserModel.FindAll(context.Background(),
		l.svcCtx.WechatUserModel.RowBuilder().Where(squirrel.Eq{"user": users}),
	)
}
