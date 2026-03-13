package svc

import (
	"sync"
	"sync/atomic"

	"wallet/service/wallet/rpc/internal/config"
)

type Wallet struct {
	Id      int64
	Balance int64
}

type ServiceContext struct {
	Config      config.Config
	Wallets     sync.Map
	WalletSeq   atomic.Int64
	WalletMu    sync.Map
	TransferSeq atomic.Int64
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
