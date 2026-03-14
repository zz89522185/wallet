package logic

import (
	"context"

	"wallet/service/wallet/api/internal/svc"
	"wallet/service/wallet/api/internal/types"
	"wallet/service/wallet/rpc/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DepositLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDepositLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DepositLogic {
	return &DepositLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DepositLogic) Deposit(req *types.DepositReq) (resp *types.DepositResp, err error) {
	rpcResp, err := l.svcCtx.WalletRpc.Deposit(l.ctx, &pb.DepositReq{
		WalletId: req.WalletId,
		Amount:   req.Amount,
	})
	if err != nil {
		return nil, err
	}

	return &types.DepositResp{
		WalletId:      rpcResp.WalletId,
		Amount:        rpcResp.Amount,
		BalanceBefore: rpcResp.BalanceBefore,
		BalanceAfter:  rpcResp.BalanceAfter,
	}, nil
}
