package script_log

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"script/script/internal/svc"
	"script/script/model"
	"time"
)

type ScriptLogRepository struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewScriptLogRepository(ctx context.Context, svcCtx *svc.ServiceContext) *ScriptLogRepository {
	return &ScriptLogRepository{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ScriptLogRepository) Insert(scriptId int64) (res sql.Result, err error) {
	ctx := context.Background()

	scriptLogModel := &model.ScriptLog{
		ScriptId:       scriptId,
		ExecutionCount: 0,
		Status:         model.ScriptLogStatusRunning,
	}

	return l.svcCtx.ScriptLogModel.Insert(ctx, scriptLogModel)
}

func (l *ScriptLogRepository) All() (pos []*model.ScriptLog, err error) {
	ctx := context.Background()

	pos, err = l.svcCtx.ScriptLogModel.FindAll(
		ctx,
		l.svcCtx.ScriptLogModel.RowBuilder().Where(squirrel.Eq{"is_delete": model.ScriptNotDeleted}).Where(squirrel.Eq{"is_enable": model.ScriptEnable}),
	)
	if nil != err {
		return
	}
	return
}

func (l *ScriptLogRepository) One(where string) (po *model.ScriptLog, err error) {
	ctx := context.Background()
	return l.svcCtx.ScriptLogModel.FindOneByQuery(ctx,
		l.svcCtx.ScriptLogModel.RowBuilder().Where(where),
	)
}

func (l *ScriptLogRepository) GetById(id int64) (po *model.ScriptLog, err error) {
	ctx := context.Background()
	return l.svcCtx.ScriptLogModel.FindOneByQuery(ctx,
		l.svcCtx.ScriptLogModel.RowBuilder().Where(squirrel.Eq{"id": id}),
	)
}

func (l *ScriptLogRepository) Update(id, status int64, result string, isEnd bool) (err error) {
	ctx := context.Background()
	scriptLogModel, err := l.GetById(id)
	if nil != err {
		return
	}
	if scriptLogModel != nil {
		nt := time.Now()
		scriptLogModel.Status = status
		scriptLogModel.Result = result
		scriptLogModel.UpdatedAt = nt
		if isEnd {
			scriptLogModel.EndAt = sql.NullTime{
				Time:  nt,
				Valid: true,
			}
		}
	}

	return l.svcCtx.ScriptLogModel.Update(ctx, scriptLogModel)

}
