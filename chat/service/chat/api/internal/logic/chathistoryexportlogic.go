package logic

import (
	"chat/common/redis"
	"chat/service/chat/api/internal/config"
	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"
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

type ChatHistoryExportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatHistoryExportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatHistoryExportLogic {
	return &ChatHistoryExportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatHistoryExportLogic) ChatHistoryExport(req *types.ChatHistoryExportReq) (resp *types.ChatHistoryExportReply, err error) {

	userId, kfId, agentId, err := l.formatCondition(req)
	if err != nil {
		return
	}
	key := l.ChatHistoryExportKey(kfId, userId, req.ChatType)
	fileUrl, err := redis.Rdb.Get(l.ctx, key).Result()
	if err == nil {
		return &types.ChatHistoryExportReply{File: fileUrl}, nil
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
	//if req.ChatType == 1 {
	//	//机器人
	//	countBuilder = countBuilder.Where(squirrel.Eq{"open_kf_id": ""})
	//	rowBuilder = rowBuilder.Where(squirrel.Eq{"open_kf_id": ""})
	//} else {
	//	countBuilder = countBuilder.Where(squirrel.Eq{"agent_id": ""})
	//	rowBuilder = rowBuilder.Where(squirrel.Eq{"agent_id": ""})
	//}

	if req.StartTime != "" {
		countBuilder = countBuilder.Where("created_at >= ?", req.StartTime)
		rowBuilder = rowBuilder.Where("created_at >= ?", req.StartTime)
	}

	if req.EndTime != "" {
		countBuilder = countBuilder.Where("created_at < ?", req.EndTime)
		rowBuilder = rowBuilder.Where("created_at < ?", req.EndTime)
	}

	chatCount, err := l.svcCtx.ChatModel.FindCount(context.Background(), countBuilder)
	if err != nil {
		return
	}
	if chatCount <= 0 {
		return nil, fmt.Errorf("没有要导出的数据")
	}
	if chatCount > MaxExportNumber {
		return nil, fmt.Errorf("导出数据太多，请加入筛选条件")
	}
	chatPo, err := l.svcCtx.ChatModel.FindAll(context.Background(), rowBuilder)
	if err != nil {
		return
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

func (l *ChatHistoryExportLogic) ChatHistoryExportKey(openKfID, userNickname string, chatType int32) string {
	return fmt.Sprintf(redis.ChatHistoryExportKey, chatType, openKfID, userNickname)
}

func (l *ChatHistoryExportLogic) GetDirName() string {
	dir := config.ChatHistoryDirectory + time.Now().Format("20060102")
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		fmt.Printf("Cannot create a file when that file already exists %v \n ", err)
	}
	return dir + "/"
}

func (l *ChatHistoryExportLogic) formatCondition(req *types.ChatHistoryExportReq) (userId, kfId string, agentId int64, err error) {
	if req.UserNickname == "" || req.KfName == "" {
		return "", "", 0, fmt.Errorf("缺少必传参数")
	}
	if req.ChatType == 2 {
		//客服聊天
		if req.UserNickname != "" {
			//get userID by UserNickName
			wechatUserPo, err := l.svcCtx.WechatUserModel.FindOneByQuery(context.Background(),
				l.svcCtx.WechatUserModel.RowBuilder().Where(squirrel.Eq{"nickname": req.UserNickname}),
			)
			if err != nil {
				return "", "", 0, err
			}

			if wechatUserPo != nil && wechatUserPo.User != "" {
				userId = wechatUserPo.User
			}
		}
		if req.KfName != "" {
			//get openKfID by OpenKfName
			kfPo, err := l.svcCtx.CustomerConfigModel.FindOneByQuery(context.Background(),
				l.svcCtx.CustomerConfigModel.RowBuilder().Where(squirrel.Eq{"kf_name": req.KfName}),
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
		userId = req.UserNickname
		if req.KfName != "" {
			//get openKfID by OpenKfName
			applicationPo, err := l.svcCtx.ApplicationConfigModel.FindOneByQuery(context.Background(),
				l.svcCtx.ApplicationConfigModel.RowBuilder().Where(squirrel.Eq{"agent_name": req.KfName}),
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
