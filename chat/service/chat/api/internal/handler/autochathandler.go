package handler

import (
	"net/http"

	"chat/service/chat/api/internal/logic"
	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AutoChatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DoGenerateActiveChatContentRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewAutoChatLogic(r.Context(), svcCtx)
		resp, err := l.AutoChat(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
