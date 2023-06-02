package logic

import (
	"chat/common/redis"
	"chat/common/util"
	"chat/service/chat/api/internal/logic/assembler"
	"chat/service/chat/api/internal/repository"
	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"
	"chat/service/chat/api/internal/vars"
	"chat/service/chat/model"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"strconv"
	"time"
)

const MaxExportNumber = 3000

type ChatRecordLogic struct {
	logx.Logger
	ctx               context.Context
	svcCtx            *svc.ServiceContext
	chatRepository    *repository.ChatRepository
	applicationConfig *repository.ApplicationConfigRepository
	customerConfig    *repository.CustomerConfigRepository
	wechatUser        *repository.WechatUserRepository
}

func NewChatRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatRecordLogic {
	return &ChatRecordLogic{
		Logger:            logx.WithContext(ctx),
		ctx:               ctx,
		svcCtx:            svcCtx,
		applicationConfig: repository.NewApplicationConfigRepository(ctx, svcCtx),
		customerConfig:    repository.NewCustomerConfigRepository(ctx, svcCtx),
		chatRepository:    repository.NewChatRepository(ctx, svcCtx),
		wechatUser:        repository.NewWechatUserRepository(ctx, svcCtx),
	}
}

func (l *ChatRecordLogic) ChatHistoryExport(req *types.GetChatListRequest) (resp *types.ChatHistoryExportReply, err error) {
	paramMD5 := util.MD5(req.Agent + req.User + req.Customer)
	key := l.ChatHistoryExportKey(paramMD5)
	fileUrl, err := redis.Rdb.Get(l.ctx, key).Result()
	if err == nil {
		return &types.ChatHistoryExportReply{File: fileUrl}, nil
	}
	req.Page = 1
	req.PageSize = MaxExportNumber
	list, err := l.GetChatList(req)
	if err != nil {
		fmt.Printf("GetSystemConfig error: %v", err)
		return
	}
	if list == nil || list.Total <= 0 || len(list.List) <= 0 {
		return nil, fmt.Errorf("没有要导出的数据")
	}

	if list.Total > MaxExportNumber {
		return nil, fmt.Errorf("导出数据太多，请加入筛选条件")
	}

	dirName := l.GetDirName()
	nowTime := time.Now().Format("150405")
	fileName := paramMD5 + nowTime + ".csv"
	fullFilePath := dirName + fileName
	fileHandle, err := os.Create(fullFilePath)
	if err != nil {
		fmt.Printf("create file error %v", err)
		return
	}
	defer fileHandle.Close()
	writer := csv.NewWriter(fileHandle)
	defer writer.Flush()

	// 写入CSV头部
	headers := []string{"ID", "user", "content", "time"}
	err = writer.Write(headers)
	if err != nil {
		return
	}

	// 写入CSV行
	i := 1
	for _, chatSon := range list.List {
		//客户聊天
		row1 := []string{strconv.Itoa(i), chatSon.User, chatSon.ReqContent, chatSon.CreatedAt}
		err = writer.Write(row1)
		if err != nil {
			return nil, err
		}
		i++
		//客服聊天
		var customer string
		if chatSon.AgentId != "" {
			customer = chatSon.AgentId
		} else {
			customer = chatSon.OpenKfId
		}
		row2 := []string{strconv.Itoa(i), customer, chatSon.ResContent, chatSon.CreatedAt}
		err = writer.Write(row2)
		if err != nil {
			return nil, err
		}

		i++
	}
	lastFile := "http://" + l.svcCtx.Config.Ip + ":" + strconv.Itoa(l.svcCtx.Config.Port) + "/api/download/chat/history?file=" + fileName
	redis.Rdb.Set(l.ctx, key, lastFile, 5*time.Minute)
	return &types.ChatHistoryExportReply{File: lastFile}, nil
}

func (l *ChatRecordLogic) ChatHistoryExportKey(key string) string {
	return fmt.Sprintf(redis.ChatHistoryExportKey, key)
}

func (l *ChatRecordLogic) GetDirName() string {
	dir := vars.ChatHistoryDirectory + time.Now().Format("20060102")
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		fmt.Printf("Cannot create a file when that file already exists %v \n ", err)
	}
	return dir + "/"
}

func (l *ChatRecordLogic) GetLastChatRecord(req *types.GetLastChatInfoRequest) (resp *types.GetLastChatInfoReply, err error) {
	chatPo, _, err := l.GetCustomerInfo(req.UserNickname, req.KfName, "", "", "created_at desc", req.ChatType, 1)
	if err != nil {
		return
	}
	if chatPo != nil && len(chatPo) > 0 && chatPo[0] != nil && chatPo[0].Id > 0 {

		return &types.GetLastChatInfoReply{
			ResContent:  chatPo[0].ReqContent,
			ReqContent:  chatPo[0].ResContent,
			CreatedTime: chatPo[0].CreatedAt.Format("2006-01-02 15:04:05"),
		}, nil

	}
	return
}

func (l *ChatRecordLogic) GetCustomerInfo(userNickname, kfName, startTime, endTime, order string, chatType, limit int32) (chatPo []*model.Chat, count int64, err error) {

	userId, kfId, agentId, err := l.FormatCondition(userNickname, kfName, chatType)
	if err != nil {
		return
	}
	countBuilder := l.svcCtx.ChatModel.CountBuilder("id")
	rowBuilder := l.svcCtx.ChatModel.RowBuilder()
	if userId != "" {
		countBuilder = countBuilder.Where(squirrel.Eq{"user": userId})
		rowBuilder = rowBuilder.Where(squirrel.Eq{"user": userId})
	}
	if kfId != "" {
		countBuilder = countBuilder.Where(squirrel.Eq{"open_kf_id": kfId})
		rowBuilder = rowBuilder.Where(squirrel.Eq{"open_kf_id": kfId})
	}
	if agentId != 0 {
		countBuilder = countBuilder.Where(squirrel.Eq{"agent_id": agentId})
		rowBuilder = rowBuilder.Where(squirrel.Eq{"agent_id": agentId})
	}

	if startTime != "" {
		countBuilder = countBuilder.Where("created_at >= ?", startTime)
		rowBuilder = rowBuilder.Where("created_at >= ?", startTime)
	}

	if endTime != "" {
		countBuilder = countBuilder.Where("created_at < ?", endTime)
		rowBuilder = rowBuilder.Where("created_at < ?", endTime)
	}

	count, err = l.svcCtx.ChatModel.FindCount(context.Background(), countBuilder)
	if err != nil {
		return
	}
	if count <= 0 {
		return nil, 0, nil
	}

	rowBuilder = rowBuilder.OrderBy(order)
	if limit != 0 {
		rowBuilder = rowBuilder.Limit(uint64(limit))
	}
	chatPo, err = l.svcCtx.ChatModel.FindAll(context.Background(), rowBuilder)
	if err != nil {
		return
	}
	return
}

func (l *ChatRecordLogic) FormatCondition(userNickname, kfName string, chatType int32) (userId, kfId string, agentId int64, err error) {
	if userNickname == "" || kfName == "" {
		return "", "", 0, fmt.Errorf("缺少必传参数")
	}
	if chatType == 2 {
		//客服聊天
		if userNickname != "" {
			//get userID by UserNickName
			wechatUserPo, err := l.svcCtx.WechatUserModel.FindOneByQuery(context.Background(),
				l.svcCtx.WechatUserModel.RowBuilder().Where(squirrel.Eq{"nickname": userNickname}),
			)
			if err != nil {
				return "", "", 0, err
			}

			if wechatUserPo != nil && wechatUserPo.User != "" {
				userId = wechatUserPo.User
			}
		}
		if kfName != "" {
			//get openKfID by OpenKfName
			kfPo, err := l.svcCtx.CustomerConfigModel.FindOneByQuery(context.Background(),
				l.svcCtx.CustomerConfigModel.RowBuilder().Where(squirrel.Eq{"kf_name": kfName}),
			)
			if err != nil {
				return "", "", 0, err
			}

			if kfPo != nil && kfPo.KfId != "" {
				kfId = kfPo.KfId
			}
		}

	} else {
		//机器人聊天
		userId = userNickname
		if kfName != "" {
			//get openKfID by OpenKfName
			applicationPo, err := l.svcCtx.ApplicationConfigModel.FindOneByQuery(context.Background(),
				l.svcCtx.ApplicationConfigModel.RowBuilder().Where(squirrel.Eq{"agent_name": kfName}),
			)
			if err != nil {
				return "", "", 0, err
			}

			if applicationPo != nil && applicationPo.AgentId != 0 {
				agentId = applicationPo.AgentId
			}
		}

	}
	return
}

func (l *ChatRecordLogic) GetChatList(req *types.GetChatListRequest) (resp *types.GetChatListPageResult, err error) {
	var agentId int64
	var user, openKfId string
	if req.Agent != "" {
		applicationPo, err := l.applicationConfig.GetByName(req.Agent)
		if nil != err {
			return nil, err
		}
		if applicationPo == nil || applicationPo.AgentId == 0 {
			return &types.GetChatListPageResult{
				List:     nil,
				Total:    0,
				Page:     req.Page,
				PageSize: req.PageSize,
			}, nil
		}
		agentId = applicationPo.AgentId
	}
	if req.ChatType == 2 {
		wechatUserPo, err := l.wechatUser.GetByName(req.User)
		if nil != err {
			return nil, err
		}
		if wechatUserPo == nil || wechatUserPo.User == "" {
			return &types.GetChatListPageResult{
				List:     nil,
				Total:    0,
				Page:     req.Page,
				PageSize: req.PageSize,
			}, nil
		}
		user = wechatUserPo.User
	} else {
		user = req.User
	}
	if req.Customer != "" {
		applicationPo, err := l.customerConfig.GetByName(req.Customer)
		if nil != err {
			return nil, err
		}
		if applicationPo == nil || applicationPo.KfId == "" {
			return &types.GetChatListPageResult{
				List:     nil,
				Total:    0,
				Page:     req.Page,
				PageSize: req.PageSize,
			}, nil
		}
		openKfId = applicationPo.KfId
	}
	chatPos, count, err := l.chatRepository.GetAll(agentId, openKfId, user, req.StartCreatedAt, req.EndCreatedAt, "id asc", uint64(req.Page), uint64(req.PageSize))
	if err != nil {
		fmt.Printf("GetSystemConfig error: %v", err)
		return
	}

	if count <= 0 || len(chatPos) <= 0 {
		return &types.GetChatListPageResult{
			List:     nil,
			Total:    0,
			Page:     req.Page,
			PageSize: req.PageSize,
		}, nil
	}
	var users, customers []string
	var applications []int64
	for _, v := range chatPos {
		if v.AgentId == 0 {
			customers = append(customers, v.OpenKfId)
			users = append(users, v.User)
		} else {
			applications = append(applications, v.AgentId)
		}
	}
	var wechatUserPos []*model.WechatUser
	var customerPos []*model.CustomerConfig
	var applicationPos []*model.ApplicationConfig
	if len(users) > 0 {

		wechatUserPos, err = l.wechatUser.GetByUsers(util.Unique(users))
		if err != nil {
			fmt.Printf("GetSystemConfig error: %v", err)
			return
		}
	}
	if len(customers) > 0 {
		customerPos, err = l.customerConfig.GetByKfIds(util.Unique(customers))
		if err != nil {
			fmt.Printf("GetSystemConfig error: %v", err)
			return
		}
	}
	if len(applications) > 0 {
		applicationPos, err = l.applicationConfig.GetByIds(util.Unique(applications))
		if err != nil {
			fmt.Printf("GetSystemConfig error: %v", err)
			return
		}
	}
	chatDtos := assembler.POTODTOGetChatList(chatPos, wechatUserPos, applicationPos, customerPos)
	return &types.GetChatListPageResult{
		List:     chatDtos,
		Total:    count,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
