package logic

import (
	"chat/common/milvus"
	"context"
	"fmt"
	"net/http"

	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCollectionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCollectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCollectionLogic {
	return &DeleteCollectionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (f *DeleteCollectionLogic) DeleteCollection(ctx context.Context, req *types.DeleteCollectionHandlerReq, r *http.Request) (resp *types.DeleteCollectionHandlerReply, err error) {
	//数据库没有
	milvusService, err := milvus.InitMilvus(f.svcCtx.Config.Embeddings.Milvus.Host, f.svcCtx.Config.Embeddings.Milvus.Username, f.svcCtx.Config.Embeddings.Milvus.Password)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = milvusService.DeleteCollection(req.CollectionName)
	return
}
