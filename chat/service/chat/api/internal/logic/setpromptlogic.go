package logic

import (
	"context"

	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetPromptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetPromptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetPromptLogic {
	return &SetPromptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetPromptLogic) SetPrompt(req *types.SetPromptReq) (resp *types.SetPromptReply, err error) {
	// todo: add your logic here and delete this line

	return
}
