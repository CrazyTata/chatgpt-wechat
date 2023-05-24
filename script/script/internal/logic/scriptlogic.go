package logic

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"os/exec"
	"script/script/internal/svc"
	"script/script/internal/types"
	"script/script/model"
	"script/script/repository/script"
	"script/script/repository/script_log"
	"script/script/util"
	"time"
)

const ScriptMaxRunTime = 2
const ErrorMaxLength = 5000

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
			l.Run(r.Id, r.Path, r.ScriptType)
		}
	}
	return &types.ScriptResponse{
		Message: "ok",
	}, nil
}

func (l *ScriptLogic) Run(scriptId int64, path, scriptType string) {
	//判断是不是有脚本正在执行的，如果有就先不处理
	where := fmt.Sprintf("script_id = %d and status = %d and created_at > '%s'", scriptId, model.ScriptLogStatusRunning, time.Now().Add(time.Duration(-ScriptMaxRunTime)*time.Minute).Format("2006-01-02 15:04:05"))
	scriptLogPo, err := l.ScriptLogRepository.One(where)
	if nil != err {
		util.Error("ScriptLogRepository.One err: " + err.Error())
		return
	}

	if scriptLogPo != nil && scriptLogPo.Id > 0 {
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
	result, err := l.HandleScript(path, scriptType)
	if nil != err {
		err = l.ScriptLogRepository.Update(logId, model.ScriptLogStatusFail, err.Error(), true)
		return
	}
	err = l.ScriptLogRepository.Update(logId, model.ScriptLogStatusSuccess, result, true)

	return
}
func (l *ScriptLogic) HandleScript(path, scriptType string) (res string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ScriptMaxRunTime)*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, scriptType, path)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	errRun := cmd.Run()
	if errRun != nil {

		originError := stderr.String()
		//errMessage := fmt.Sprint(originError)
		var newError string
		if len(originError) > ErrorMaxLength {
			newError = originError[:ErrorMaxLength]
		} else {
			newError = originError
		}
		err = errors.New(errRun.Error() + newError)

		util.Info("HandleScript " + errRun.Error() + originError)
		return
	}
	// cmd.Run()执行成功，输出正常信息
	res = stdout.String()
	return
}
