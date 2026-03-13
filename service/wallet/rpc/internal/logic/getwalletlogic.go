package logic

import (
	"context"
	"sync"

	"wallet/service/wallet/rpc/internal/svc"
	"wallet/service/wallet/rpc/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetWalletLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWalletLogic {
	return &GetWalletLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetWalletLogic) GetWallet(in *pb.GetWalletReq) (*pb.GetWalletResp, error) {
	mu, _ := l.svcCtx.WalletMu.LoadOrStore(in.WalletId, &sync.Mutex{})
	mu.(*sync.Mutex).Lock()
	defer mu.(*sync.Mutex).Unlock()

	val, ok := l.svcCtx.Wallets.Load(in.WalletId)
	if !ok {
		return nil, status.Error(codes.NotFound, "wallet not found")
	}

	w := val.(*svc.Wallet)
	return &pb.GetWalletResp{
		WalletId: w.Id,
		Balance:  w.Balance,
	}, nil
}
