package script

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"script/script/internal/svc"
	"script/script/model"
	"time"
)

type ScriptRepository struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewScriptRepository(ctx context.Context, svcCtx *svc.ServiceContext) *ScriptRepository {
	return &ScriptRepository{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ScriptRepository) InsertScript(name, path, scriptType string) (err error) {
	ctx := context.Background()
	if name == "" {
		return
	}
	//先查询下
	scriptPo, err := l.svcCtx.ScriptModel.FindOneByQuery(ctx,
		l.svcCtx.ScriptModel.RowBuilder().Where(squirrel.Eq{"name": name}),
	)
	if err != nil {
		fmt.Printf("InsertWechatUser FindOneByQuery error: %v", err)
		return
	}

	if scriptPo != nil && scriptPo.Id > 0 {
		scriptPo.Name = name
		scriptPo.Path = path
		scriptPo.IsDelete = model.ScriptNotDeleted
		scriptPo.IsEnable = model.ScriptEnable
		scriptPo.UpdatedAt = time.Now()
		err = l.svcCtx.ScriptModel.Update(ctx, scriptPo)
		return
	}
	_, err = l.svcCtx.ScriptModel.Insert(ctx, &model.Script{
		Name:       name,
		Path:       path,
		IsDelete:   model.ScriptNotDeleted,
		IsEnable:   model.ScriptEnable,
		ScriptType: scriptType,
	})

	return
}

func (l *ScriptRepository) All() (pos []*model.Script, err error) {
	ctx := context.Background()

	pos, err = l.svcCtx.ScriptModel.FindAll(
		ctx,
		l.svcCtx.ScriptModel.RowBuilder().Where(squirrel.Eq{"is_delete": model.ScriptNotDeleted}).Where(squirrel.Eq{"is_enable": model.ScriptEnable}),
	)
	if nil != err {
		return
	}
	return
}
