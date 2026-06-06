package service

import (
	"context"
	"time"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/google/uuid"
	"github.com/mkfolder/moxie/internal/db"
	"github.com/mkfolder/moxie/internal/models"
	"github.com/mkfolder/moxie/pkg/http"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const memoProgramID = "MemoSq4gqABAXKb96qnH8TysNcWxMyWCqXgDLGmfcHr"

type AuthConfig struct {
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type Service struct {
	log       *zap.SugaredLogger
	db        *gorm.DB
	http      *http.Client
	rpc       *rpc.Client
	orders    db.Repository[models.Order]
	merchants db.Repository[models.Merchant]
	authCfg   AuthConfig
}

type NewServiceParams struct {
	Log    *zap.SugaredLogger
	HTTP   *http.Client
	GormDB *gorm.DB
	RPC    string
	Auth   AuthConfig
}

func New(params NewServiceParams) *Service {
	if params.Log == nil {
		panic("invalid service.New arguments: logger is nil")
	}

	if params.GormDB == nil {
		panic("invalid service.New arguments: gorm db is nil")
	}

	rpcClient := rpc.New(rpc.DevNet_RPC)
	ordersRepository := db.NewGormRepository[models.Order](params.GormDB)
	merchantsRepository := db.NewGormRepository[models.Merchant](params.GormDB)
	return &Service{
		log:       params.Log,
		http:      params.HTTP,
		db:        params.GormDB,
		orders:    ordersRepository,
		merchants: merchantsRepository,
		rpc:       rpcClient,
		authCfg:   params.Auth,
	}
}

func (s *Service) FindMerchantByID(ctx context.Context, id uuid.UUID) (*models.Merchant, error) {
	return s.merchants.Find(ctx, id)
}

func (s *Service) PingDB(ctx context.Context) error {
	return s.db.WithContext(ctx).Raw("SELECT 1").Error
}
