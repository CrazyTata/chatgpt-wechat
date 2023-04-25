package logic

import (
	"chat/common/xerr"
	"chat/service/chat/model"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"

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

	//先查询下
	promptPo, err := l.svcCtx.CustomerPromptModel.FindOneByQuery(context.Background(),
		l.svcCtx.CustomerPromptModel.RowBuilder().Where(squirrel.Eq{"kf_id": req.KfId}),
	)
	if err != nil {
		fmt.Printf("ListPrompt error: %v", err)
		return
	}
	if promptPo == nil {
		err = fmt.Errorf("当前客服id不存在")
		return
	}
	if err := l.svcCtx.CustomerPromptModel.Update(l.ctx, &model.CustomerPrompt{
		Id:     promptPo.Id,
		KfName: promptPo.KfName,
		KfId:   promptPo.KfId,
		Prompt: req.Prompt,
	}); err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("更新失败"), "更新失败,error: %v", err)
	}
	return &types.SetPromptReply{
		Message: "ok",
	}, nil
}
