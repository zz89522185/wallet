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

type DepositLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDepositLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DepositLogic {
	return &DepositLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DepositLogic) Deposit(in *pb.DepositReq) (*pb.DepositResp, error) {
	if in.Amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "deposit amount must be positive")
	}

	_, ok := l.svcCtx.Wallets.Load(in.WalletId)
	if !ok {
		return nil, status.Error(codes.NotFound, "wallet not found")
	}

	mu, _ := l.svcCtx.WalletMu.LoadOrStore(in.WalletId, &sync.Mutex{})
	mu.(*sync.Mutex).Lock()
	defer mu.(*sync.Mutex).Unlock()

	val, _ := l.svcCtx.Wallets.Load(in.WalletId)
	w := val.(*svc.Wallet)

	balanceBefore := w.Balance
	w.Balance += in.Amount

	return &pb.DepositResp{
		WalletId:      in.WalletId,
		Amount:        in.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  w.Balance,
	}, nil
}
