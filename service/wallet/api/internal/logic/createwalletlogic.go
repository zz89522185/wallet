package logic

import (
	"context"

	"wallet/service/wallet/api/internal/svc"
	"wallet/service/wallet/api/internal/types"
	"wallet/service/wallet/rpc/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateWalletLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateWalletLogic {
	return &CreateWalletLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateWalletLogic) CreateWallet(req *types.CreateWalletReq) (resp *types.CreateWalletResp, err error) {
	rpcResp, err := l.svcCtx.WalletRpc.CreateWallet(l.ctx, &pb.CreateWalletReq{
		UserId:   req.UserId,
		Currency: req.Currency,
	})
	if err != nil {
		return nil, err
	}

	return &types.CreateWalletResp{
		WalletId: rpcResp.WalletId,
		Balance:  rpcResp.Balance,
	}, nil
}
