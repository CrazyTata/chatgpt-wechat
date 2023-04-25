package handler

import (
	"chat/common/validator"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"chat/common/response"
	"chat/service/chat/api/internal/logic"
	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"
)

func DeleteCollectionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteCollectionHandlerReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamError(r, w, err)
			return
		}
		//validate check
		if err := validator.Validate.StructCtx(r.Context(), req); err != nil {
			errMap := validator.Translate(err, &req)
			for _, errFormat := range errMap {
				response.ParamError(r, w, fmt.Errorf(errFormat))
				return
			}
			response.ParamError(r, w, err)
			return
		}

		l := logic.NewDeleteCollectionLogic(r.Context(), svcCtx)
		resp, err := l.DeleteCollection(r.Context(), &req, r)
		response.Response(r, w, resp, err)
	}
}
