package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"wallet/service/wallet/api/internal/logic"
	"wallet/service/wallet/api/internal/svc"
	"wallet/service/wallet/api/internal/types"
)

func GetWalletHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetWalletReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetWalletLogic(r.Context(), svcCtx)
		resp, err := l.GetWallet(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
