package logic

import (
	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"
	"context"
	"fmt"
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
	promptPo, err := l.svcCtx.CustomerPromptModel.FindAll(context.Background(),
		l.svcCtx.CustomerPromptModel.RowBuilder(),
	)
	if err != nil {
		fmt.Printf("ListPrompt error: %v", err)
		return
	}
	var p types.ListPromptReply
	if len(promptPo) > 0 {
		for _, v := range promptPo {
			p.List = append(p.List, types.ListPromptReplyData{
				Id:     v.Id,
				KfId:   v.KfId,
				Prompt: v.Prompt,
			})
		}
	}
	resp = &p
	return
}
