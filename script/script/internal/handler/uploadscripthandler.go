package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"script/script/internal/logic"
	"script/script/internal/svc"
	"script/script/internal/types"
)

func UploadScriptHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadScriptRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUploadScriptLogic(r.Context(), svcCtx)
		resp, err := l.UploadScript(&req, r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
