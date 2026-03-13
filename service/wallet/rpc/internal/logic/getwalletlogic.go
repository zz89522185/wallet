package logic

import (
	"context"

	"wallet/service/wallet/rpc/internal/svc"
	"wallet/service/wallet/rpc/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
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
	// todo: add your logic here and delete this line

	return &pb.GetWalletResp{}, nil
}
