package logic

import (
	"context"

	"wallet/service/wallet/rpc/internal/svc"
	"wallet/service/wallet/rpc/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateWalletLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateWalletLogic {
	return &CreateWalletLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateWalletLogic) CreateWallet(in *pb.CreateWalletReq) (*pb.CreateWalletResp, error) {
	id := l.svcCtx.WalletSeq.Add(1)
	w := &svc.Wallet{
		Id:      id,
		Balance: 0,
	}
	l.svcCtx.Wallets.Store(id, w)

	return &pb.CreateWalletResp{
		WalletId: id,
		Balance:  0,
	}, nil
}
