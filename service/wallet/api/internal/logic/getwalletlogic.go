package logic

import (
	"context"

	"wallet/service/wallet/api/internal/svc"
	"wallet/service/wallet/api/internal/types"
	"wallet/service/wallet/rpc/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWalletLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWalletLogic {
	return &GetWalletLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWalletLogic) GetWallet(req *types.GetWalletReq) (resp *types.GetWalletResp, err error) {
	rpcResp, err := l.svcCtx.WalletRpc.GetWallet(l.ctx, &pb.GetWalletReq{
		WalletId: req.WalletId,
	})
	if err != nil {
		return nil, err
	}

	return &types.GetWalletResp{
		WalletId: rpcResp.WalletId,
		Balance:  rpcResp.Balance,
		Currency: rpcResp.Currency,
	}, nil
}
