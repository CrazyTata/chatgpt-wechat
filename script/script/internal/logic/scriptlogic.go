package logic

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"script/script/model"
	"script/script/repository/script"
	"script/script/repository/script_log"
	"script/script/util"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"script/script/internal/svc"
	"script/script/internal/types"
)

const ScriptMaxRunTime = 10

type ScriptLogic struct {
	logx.Logger
	ctx                 context.Context
	svcCtx              *svc.ServiceContext
	ScriptRepository    *script.ScriptRepository
	ScriptLogRepository *script_log.ScriptLogRepository
}

func NewScriptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ScriptLogic {
	return &ScriptLogic{
		Logger:              logx.WithContext(ctx),
		ctx:                 ctx,
		svcCtx:              svcCtx,
		ScriptRepository:    script.NewScriptRepository(ctx, svcCtx),
		ScriptLogRepository: script_log.NewScriptLogRepository(ctx, svcCtx),
	}
}

func (l *ScriptLogic) Script(req *types.ScriptRequest) (resp *types.ScriptResponse, err error) {
	// 获取所有的可用脚本
	scriptPos, err := l.ScriptRepository.All()
	if nil != err {
		util.Error("ScriptRepository.All err: " + err.Error())
		return
	}

	if len(scriptPos) > 0 {
		for _, r := range scriptPos {
			r := r
			l.Run(r.Id, r.Path)
		}
	}
	return &types.ScriptResponse{
		Message: "ok",
	}, nil
}

func (l *ScriptLogic) Run(scriptId int64, path string) {
	//判断是不是有脚本正在执行的，如果有就先不处理
	where := fmt.Sprintf("script_id = %d and status = %d and created_at > '%s'", scriptId, model.ScriptLogStatusRunning, time.Now().Add(time.Duration(-ScriptMaxRunTime)*time.Minute).Format("2006-01-02 15:04:05"))
	scriptLogPo, err := l.ScriptLogRepository.One(where)
	if nil != err {
		fmt.Println("ScriptLogRepository.One err: " + err.Error())
		util.Error("ScriptLogRepository.One err: " + err.Error())
		return
	}

	if scriptLogPo != nil && scriptLogPo.Id > 0 {
		fmt.Println(fmt.Sprintf("the task is running scriptId = %d  scriptLogId = %d ", scriptId, scriptLogPo.Id))
		util.Info(fmt.Sprintf("the task is running scriptId = %d  scriptLogId = %d ", scriptId, scriptLogPo.Id))
		return
	}
	res, err := l.ScriptLogRepository.Insert(scriptId)
	if nil != err {
		util.Error("ScriptLogRepository.Insert err: " + err.Error())
		return
	}
	logId, err := res.LastInsertId()

	if nil != err {
		util.Error("ScriptLogRepository.LastInsertId err: " + err.Error())
		return
	}
	result, err := l.HandleScript(path)
	if nil != err {
		err = l.ScriptLogRepository.Update(logId, model.ScriptLogStatusFail, err.Error(), true)
		return
	}
	err = l.ScriptLogRepository.Update(logId, model.ScriptLogStatusSuccess, result, true)

	return
}
func (l *ScriptLogic) HandleScript(path string) (res string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ScriptMaxRunTime)*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "python3", path)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	errRun := cmd.Run()
	if errRun != nil {
		errMessage := fmt.Sprint(errRun) + ": " + stderr.String()
		err = errors.New(errMessage)
		util.Info("HandleScript" + errMessage)
		return
	}
	// cmd.Run()执行成功，输出正常信息
	res = stdout.String()
	return
}
