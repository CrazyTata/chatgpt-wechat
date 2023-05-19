package logic

import (
	"context"

	"script/script/internal/svc"
	"script/script/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ScriptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewScriptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ScriptLogic {
	return &ScriptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ScriptLogic) Script(req *types.ScriptRequest) (resp *types.ScriptResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
