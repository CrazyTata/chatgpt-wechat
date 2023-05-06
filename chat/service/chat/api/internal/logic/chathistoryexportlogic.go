package logic

import (
	"chat/common/redis"
	"chat/service/chat/api/internal/config"
	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
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
	key := l.ChatHistoryExportKey(req.OpenKfID, req.UserNickname)
	fileUrl, err := redis.Rdb.Get(l.ctx, key).Result()
	if err == nil {
		return &types.ChatHistoryExportReply{File: fileUrl}, nil
	}
	//get userID by UserNickName
	wechatUserPo, err := l.svcCtx.WechatUserModel.FindOneByQuery(context.Background(),
		l.svcCtx.WechatUserModel.RowBuilder().Where(squirrel.Eq{"nickname": req.UserNickname}),
	)
	if err != nil {
		return nil, err
	}

	if wechatUserPo == nil || wechatUserPo.User == "" {
		return nil, fmt.Errorf("当前用户下没有要导出的数据")
	}

	chatCount, err := l.svcCtx.ChatModel.FindCount(context.Background(),
		l.svcCtx.ChatModel.CountBuilder("id").Where(squirrel.Eq{"user": wechatUserPo.User}).Where(squirrel.Eq{"open_kf_id": req.OpenKfID}),
	)
	if err != nil {
		return
	}
	if chatCount <= 0 {
		return nil, fmt.Errorf("没有要导出的数据")
	}
	if chatCount > MaxExportNumber {
		return nil, fmt.Errorf("导出数据太多，请加入时间筛选")
	}
	chatPo, err := l.svcCtx.ChatModel.FindAll(context.Background(),
		l.svcCtx.ChatModel.RowBuilder().Where(squirrel.Eq{"user": wechatUserPo.User}).Where(squirrel.Eq{"open_kf_id": req.OpenKfID}),
	)
	if err != nil {
		return
	}
	dirName := l.GetDirName()
	nowTime := time.Now().Format("150405")

	fileName := wechatUserPo.User + nowTime + ".csv"
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
	headers := []string{"id", "问题", "答案"}
	err = writer.Write(headers)
	if err != nil {
		return
	}

	// 写入CSV行
	for _, chatSon := range chatPo {
		row := []string{strconv.Itoa(int(chatSon.Id)), chatSon.ReqContent, chatSon.ResContent}
		err = writer.Write(row)
		if err != nil {
			return nil, err
		}
	}

	return &types.ChatHistoryExportReply{File: fileName}, nil
}

func (l *ChatHistoryExportLogic) ChatHistoryExportKey(openKfID, userNickname string) string {
	return fmt.Sprintf(redis.ChatHistoryExportKey, openKfID, userNickname)
}

func (l *ChatHistoryExportLogic) GetDirName() string {
	dir := config.ChatHistoryDirectory + time.Now().Format("20060102")
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		fmt.Printf("Cannot create a file when that file already exists %v \n ", err)
	}
	return dir + "/"
}

func handleDownload(w http.ResponseWriter, filePath string) {
	filepath := "/path/to/file"
	file, err := os.Open(filepath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+filepath)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	io.Copy(w, file)
}
