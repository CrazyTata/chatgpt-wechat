package logic

import (
	"context"

	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPromptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPromptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPromptLogic {
	return &ListPromptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPromptLogic) ListPrompt(req *types.ListPromptReq) (resp *types.ListPromptReply, err error) {
	// todo: add your logic here and delete this line

	return
}
