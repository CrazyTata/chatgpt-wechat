package logic

import (
	"chat/common/wecom"
	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"
	"chat/service/chat/model"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/logx"
)

const SyncNumber = 90

type SyncWechatUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSyncWechatUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncWechatUserLogic {
	return &SyncWechatUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SyncWechatUserLogic) SyncWechatUser(req *types.SyncWechatUserReq) (resp *types.SyncWechatUserReply, err error) {

	count, err := l.svcCtx.WechatUserModel.FindCount(context.Background(),
		l.svcCtx.WechatUserModel.CountBuilder("id").Where(squirrel.Eq{"nickname": ""}),
	)
	if err != nil {
		return
	}
	if count <= 0 {
		return nil, fmt.Errorf("没有要同步的数据")
	}
	//workNumber := math.Ceil(float64(count) / SyncNumber)

	l.HandleSyncWechatUser(0)

	return
}

func (l *SyncWechatUserLogic) HandleSyncWechatUser(page int64) {
	limit := uint64(SyncNumber)
	offset := uint64(page * SyncNumber)
	pos, err := l.svcCtx.WechatUserModel.FindAll(context.Background(),
		l.svcCtx.WechatUserModel.RowBuilder().Where(squirrel.Eq{"nickname": ""}).Offset(offset).Limit(limit),
	)
	if err != nil {
		fmt.Printf("HandleSyncWechatUser error: %+v", err)
		return
	}
	if len(pos) > 0 {
		var users []string
		posItems := make(map[string]int64)

		for _, v := range pos {
			users = append(users, v.User)
			posItems[v.User] = v.Id
		}
		originWechatInfo, err := wecom.GetCustomer(users)
		if err != nil {
			fmt.Printf("GetCustomer error: %+v", err)
			return
		}

		for _, wechat := range originWechatInfo {
			if wechat.Nickname != "" {
				err = l.svcCtx.WechatUserModel.Update(context.Background(), &model.WechatUser{
					Id:       posItems[wechat.ExternalUserID],
					User:     wechat.ExternalUserID,
					Nickname: wechat.Nickname,
					Avatar:   wechat.Avatar,
					Gender:   wechat.Gender,
					Unionid:  wechat.Unionid,
				})
				if err != nil {
					return
				}
			}
		}
	}
	return
}
