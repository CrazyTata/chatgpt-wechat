package logic

import (
	"context"
	"fmt"

	"cron/cron/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CronLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCronLogic(svcCtx *svc.ServiceContext) *CronLogic {
	ctx := context.Background()
	return &CronLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CronLogic) PythonScript() {
	// todo: add your logic here and delete this line
	fmt.Print("test tata")
	return
}
