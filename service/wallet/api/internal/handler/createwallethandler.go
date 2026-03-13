package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"wallet/service/wallet/api/internal/logic"
	"wallet/service/wallet/api/internal/svc"
	"wallet/service/wallet/api/internal/types"
)

func CreateWalletHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateWalletReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCreateWalletLogic(r.Context(), svcCtx)
		resp, err := l.CreateWallet(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
