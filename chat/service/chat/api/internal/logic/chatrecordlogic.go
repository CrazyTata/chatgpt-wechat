package logic

import (
	"chat/common/redis"
	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"
	"chat/service/chat/api/internal/vars"
	"chat/service/chat/model"
	"context"
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"strconv"
	"time"
)

const MaxExportNumber = 5000

type ChatRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatRecordLogic {
	return &ChatRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatRecordLogic) ChatHistoryExport(req *types.ChatHistoryExportReq) (resp *types.ChatHistoryExportReply, err error) {

	userId, kfId, _, err := l.FormatCondition(req.UserNickname, req.KfName, req.ChatType)
	if err != nil {
		return
	}
	key := l.ChatHistoryExportKey(kfId, userId, req.ChatType)
	fileUrl, err := redis.Rdb.Get(l.ctx, key).Result()
	if err == nil {
		return &types.ChatHistoryExportReply{File: fileUrl}, nil
	}
	chatPo, chatCount, err := l.GetCustomerInfo(req.UserNickname, req.KfName, req.StartTime, req.EndTime, "id asc", req.ChatType, MaxExportNumber)
	if chatCount <= 0 {
		return nil, fmt.Errorf("没有要导出的数据")
	}
	if chatCount > MaxExportNumber {
		return nil, fmt.Errorf("导出数据太多，请加入筛选条件")
	}

	dirName := l.GetDirName()
	nowTime := time.Now().Format("150405")
	str := userId + kfId + strconv.Itoa(int(req.ChatType))
	hash := md5.Sum([]byte(str))
	fileName := hex.EncodeToString(hash[:]) + nowTime + ".csv"
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
	for _, chatSon := range chatPo {
		//客户聊天
		row1 := []string{strconv.Itoa(i), req.UserNickname, chatSon.ReqContent, chatSon.CreatedAt.Format("2006-01-02 15:04:05")}
		err = writer.Write(row1)
		if err != nil {
			return nil, err
		}
		i++
		//客服聊天
		row2 := []string{strconv.Itoa(i), req.KfName, chatSon.ResContent, chatSon.CreatedAt.Format("2006-01-02 15:04:05")}
		err = writer.Write(row2)
		if err != nil {
			return nil, err
		}

		i++
	}
	lastFile := l.svcCtx.Config.Ip + ":" + strconv.Itoa(l.svcCtx.Config.Port) + "/api/download/chat/history?file=" + fileName
	redis.Rdb.Set(l.ctx, key, lastFile, 5*time.Minute)
	return &types.ChatHistoryExportReply{File: lastFile}, nil
}

func (l *ChatRecordLogic) ChatHistoryExportKey(openKfID, userNickname string, chatType int32) string {
	return fmt.Sprintf(redis.ChatHistoryExportKey, chatType, openKfID, userNickname)
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
