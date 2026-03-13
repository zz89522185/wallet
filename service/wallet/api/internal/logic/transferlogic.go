package logic

import (
	"context"

	"wallet/service/wallet/api/internal/svc"
	"wallet/service/wallet/api/internal/types"
	"wallet/service/wallet/rpc/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type TransferLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTransferLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransferLogic {
	return &TransferLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TransferLogic) Transfer(req *types.TransferReq) (resp *types.TransferResp, err error) {
	rpcResp, err := l.svcCtx.WalletRpc.Transfer(l.ctx, &pb.TransferReq{
		FromWalletId: req.FromWalletId,
		ToWalletId:   req.ToWalletId,
		Amount:       req.Amount,
		Remark:       req.Remark,
	})
	if err != nil {
		return nil, err
	}

	return &types.TransferResp{
		TransferId: rpcResp.TransferId,
	}, nil
}
