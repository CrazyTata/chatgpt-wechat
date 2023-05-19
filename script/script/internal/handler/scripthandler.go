package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"script/script/internal/logic"
	"script/script/internal/svc"
	"script/script/internal/types"
)

func ScriptHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ScriptRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewScriptLogic(r.Context(), svcCtx)
		resp, err := l.Script(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
