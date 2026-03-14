package logic

import (
	"context"
	"fmt"
	"sync"

	"wallet/service/wallet/rpc/internal/svc"
	"wallet/service/wallet/rpc/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TransferLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTransferLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransferLogic {
	return &TransferLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TransferLogic) Transfer(in *pb.TransferReq) (*pb.TransferResp, error) {
	// 2.1 参数校验
	if in.FromWalletId == in.ToWalletId {
		return nil, status.Error(codes.InvalidArgument, "cannot transfer to the same wallet")
	}
	if in.Amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "transfer amount must be positive")
	}

	// 2.2 按 walletId 升序确定加锁顺序，防止死锁
	firstId, secondId := in.FromWalletId, in.ToWalletId
	if firstId > secondId {
		firstId, secondId = secondId, firstId
	}

	firstMu := l.getWalletMu(firstId)
	secondMu := l.getWalletMu(secondId)

	firstMu.Lock()
	defer firstMu.Unlock()
	secondMu.Lock()
	defer secondMu.Unlock()

	// 2.3 钱包查找
	fromVal, ok := l.svcCtx.Wallets.Load(in.FromWalletId)
	if !ok {
		return nil, status.Error(codes.NotFound, "source wallet not found")
	}
	toVal, ok := l.svcCtx.Wallets.Load(in.ToWalletId)
	if !ok {
		return nil, status.Error(codes.NotFound, "target wallet not found")
	}

	fromWallet := fromVal.(*svc.Wallet)
	toWallet := toVal.(*svc.Wallet)

	// 2.4 余额检查与原子转账
	if fromWallet.Balance < in.Amount {
		return nil, status.Error(codes.FailedPrecondition, "insufficient balance")
	}

	fromWallet.Balance -= in.Amount
	toWallet.Balance += in.Amount

	// 2.5 生成 transferId
	seq := l.svcCtx.TransferSeq.Add(1)
	transferId := fmt.Sprintf("TX-%d", seq)

	return &pb.TransferResp{
		TransferId: transferId,
	}, nil
}

func (l *TransferLogic) getWalletMu(walletId int64) *sync.Mutex {
	mu, _ := l.svcCtx.WalletMu.LoadOrStore(walletId, &sync.Mutex{})
	return mu.(*sync.Mutex)
}
