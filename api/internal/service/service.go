package service

import (
	"context"

	"github.com/Makefolder/moxie/internal/db"
	"github.com/Makefolder/moxie/internal/helius"
	"github.com/Makefolder/moxie/internal/models"
	"github.com/Makefolder/moxie/pkg/http"
	"github.com/gagliardetto/solana-go/rpc"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const memoProgramID = "MemoSq4gqABAXKb96qnH8TysNcWxMyWCqXgDLGmfcHr"

type Service struct {
	log       *zap.SugaredLogger
	db        *gorm.DB
	hc        *helius.HeliusClient
	http      *http.Client
	rpc       *rpc.Client
	orders    db.Repository[models.Order]
	merchants db.Repository[models.Merchant]
}

type NewServiceParams struct {
	Log    *zap.SugaredLogger
	HC     *helius.HeliusClient
	HTTP   *http.Client
	GormDB *gorm.DB
	RPC    string
}

func New(params NewServiceParams) *Service {
	if params.Log == nil {
		panic("invalid service.New arguments: logger is nil")
	}

	if params.HC == nil {
		panic("invalid service.New arguments: helius client is nil")
	}

	if params.GormDB == nil {
		panic("invalid service.New arguments: gorm db is nil")
	}

	rpcClient := rpc.New(rpc.DevNet_RPC)
	ordersRepository := db.NewGormRepository[models.Order](params.GormDB)
	merchantsRepository := db.NewGormRepository[models.Merchant](params.GormDB)
	return &Service{
		log:       params.Log,
		hc:        params.HC,
		http:      params.HTTP,
		db:        params.GormDB,
		orders:    ordersRepository,
		merchants: merchantsRepository,
		rpc:       rpcClient,
	}
}

func (s *Service) PingDB(ctx context.Context) error {
	return s.db.WithContext(ctx).Raw("SELECT 1").Error
}
