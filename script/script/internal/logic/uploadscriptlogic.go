package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"os"
	"script/script/internal/svc"
	"script/script/internal/types"
	"script/script/internal/vars"
	"script/script/repository/script"
	"script/script/util"
)

type UploadScriptLogic struct {
	logx.Logger
	ctx              context.Context
	svcCtx           *svc.ServiceContext
	ScriptRepository *script.ScriptRepository
}

func NewUploadScriptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadScriptLogic {
	return &UploadScriptLogic{
		Logger:           logx.WithContext(ctx),
		ctx:              ctx,
		svcCtx:           svcCtx,
		ScriptRepository: script.NewScriptRepository(ctx, svcCtx),
	}
}

func (l *UploadScriptLogic) UploadScript(req *types.UploadScriptRequest, r *http.Request) (resp *types.UploadScriptResponse, err error) {
	// 1. parse input , type multipart/form-data
	r.ParseMultipartForm(1000)
	// 获取name字段值
	name := r.FormValue("name")
	scriptType := r.FormValue("script_type")
	if name == "" {
		err = errors.New("file name not exist")
		return
	}
	if scriptType == "" {
		err = errors.New("script type not exist")
		return
	}
	AllowScriptType := []string{"python3"}

	if !util.InArray(AllowScriptType, scriptType) {
		err = errors.New("script type is not allow")
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		util.Error("Error retrieving file from form-data err:" + err.Error())
		return
	}
	defer file.Close()
	//if !strings.Contains(handler.Filename, "py") {
	//	err = errors.New("please upload python script")
	//	return
	//}

	fileName := vars.ScriptDir + handler.Filename
	_, err = util.RenameFileIfExists(fileName)
	if err != nil {
		util.Error("Error retrieving file from form-data err:" + err.Error())
		return
	}
	dataBuf := make([]byte, handler.Size)
	_, err = file.Read(dataBuf)
	if err != nil {
		util.Error("read file error err:" + err.Error())
		return
	}
	err = os.MkdirAll(vars.ScriptDir, 0755)
	if err != nil {
		util.Info("Cannot create a file when that file already exists err:" + err.Error())
	}

	fileHandle, err := os.Create(fileName)
	if err != nil {
		util.Error("create file error err:" + err.Error())
		return
	}
	_, err = fileHandle.Write(dataBuf)
	if err != nil {
		util.Error("write file error err:" + err.Error())
		return
	}
	util.Info("upload file success")

	//save data
	// 再去插入数据
	err = l.ScriptRepository.InsertScript(name, fileName, scriptType)
	if err != nil {
		util.Error("insert file error err:" + err.Error())
		return
	}
	util.Info("save file success")
	return &types.UploadScriptResponse{
		Message: "ok",
	}, nil
}
