package logic

import (
	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"
	"chat/service/chat/api/internal/vars"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type DownloadChatHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDownloadChatHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadChatHistoryLogic {
	return &DownloadChatHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DownloadChatHistoryLogic) DownloadChatHistory(req *types.DownloadChatHistoryReq, w http.ResponseWriter, r *http.Request) {
	filePath := vars.ChatHistoryDirectory + time.Now().Format("20060102") + "/" + req.File

	file, err := os.Open(filePath)
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

	w.Header().Set("Content-Disposition", "attachment; filename="+filePath)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	http.ServeContent(w, r, filePath, fileInfo.ModTime(), file)
	io.Copy(w, file)

}
