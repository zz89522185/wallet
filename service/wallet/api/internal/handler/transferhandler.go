package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"wallet/service/wallet/api/internal/logic"
	"wallet/service/wallet/api/internal/svc"
	"wallet/service/wallet/api/internal/types"
)

func TransferHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TransferReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewTransferLogic(r.Context(), svcCtx)
		resp, err := l.Transfer(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
