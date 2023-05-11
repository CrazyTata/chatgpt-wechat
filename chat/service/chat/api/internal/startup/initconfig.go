package startup

import (
	"chat/common/wecom"
	"chat/service/chat/api/internal/svc"
	"context"
	"fmt"
)

func InitConfig(svcCtx *svc.ServiceContext) (result []wecom.Application) {

	//get config
	applicationConfigPo, err := svcCtx.ApplicationConfigModel.FindAll(context.Background(),
		svcCtx.ApplicationConfigModel.RowBuilder(),
	)

	if err != nil {
		fmt.Printf("InitConfig:%+v", err)
		return
	}
	if applicationConfigPo != nil {
		for _, v := range applicationConfigPo {
			result = append(result, wecom.Application{
				AgentID:     v.AgentId,
				AgentSecret: v.AgentSecret,
			})
			if v.GroupEnable {
				fmt.Println("初始化群聊", v.GroupName, v.GroupChatId, v.AgentSecret, v.AgentId)
				go wecom.InitGroup(v.GroupName, v.GroupChatId, v.AgentSecret, v.AgentId)
			}

		}

	}
	return
}
