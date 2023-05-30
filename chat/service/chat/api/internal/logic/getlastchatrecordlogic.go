package logic

import (
	"context"

	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLastChatRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLastChatRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLastChatRecordLogic {
	return &GetLastChatRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLastChatRecordLogic) GetLastChatRecord(req *types.GetLastChatInfoRequest) (resp *types.GetLastChatInfoReply, err error) {
	// todo: add your logic here and delete this line

	return
}
