package svc

import (
	"wallet/service/wallet/api/internal/config"
	"wallet/service/wallet/rpc/wallet"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	WalletRpc wallet.Wallet
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		WalletRpc: wallet.NewWallet(zrpc.MustNewClient(c.WalletRpc)),
	}
}
