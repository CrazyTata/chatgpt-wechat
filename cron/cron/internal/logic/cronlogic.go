package logic

import (
	"context"
	"cron/cron/internal/client"
	"cron/cron/internal/svc"
	"cron/cron/util"
	"github.com/zeromicro/go-zero/core/logx"
)

type CronLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	client *client.HttpService
}

func NewCronLogic(svcCtx *svc.ServiceContext) *CronLogic {
	ctx := context.Background()
	return &CronLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		client: client.NewHttpService(),
	}
}

func (l *CronLogic) RunScript() {
	util.Info("start to deal RunScript")
	l.client.RunScript()
	return
}
