package service

import (
	"context"
	"time"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/google/uuid"
	"github.com/mkfolder/moxie/internal/db"
	"github.com/mkfolder/moxie/internal/models"
	"github.com/mkfolder/moxie/pkg/http"
	"github.com/mkfolder/moxie/pkg/s3client"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthConfig struct {
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	EncryptionKey   string
}

type Service struct {
	log       *zap.SugaredLogger
	db        *gorm.DB
	http      *http.Client
	rpc       *rpc.Client
	orders    db.Repository[models.Order]
	merchants db.Repository[models.Merchant]
	authCfg   AuthConfig
	s3        *s3client.Client
}

type NewServiceParams struct {
	Log    *zap.SugaredLogger
	HTTP   *http.Client
	GormDB *gorm.DB
	RPC    string
	Auth   AuthConfig
	S3     *s3client.Client
}

func New(params NewServiceParams) *Service {
	if params.Log == nil {
		panic("invalid service.New arguments: logger is nil")
	}

	if params.GormDB == nil {
		panic("invalid service.New arguments: gorm db is nil")
	}

	if params.HTTP == nil {
		panic("invalid service.New arguments: http client is nil")
	}

	if params.S3 == nil {
		panic("invalid service.New arguments: s3 client is nil")
	}

	if params.RPC == "" {
		params.RPC = rpc.DevNet_RPC
	}

	rpcClient := rpc.New(params.RPC)
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
		s3:        params.S3,
	}
}

func (s *Service) FindMerchantByID(ctx context.Context, id uuid.UUID) (*models.Merchant, error) {
	return s.merchants.Find(ctx, id)
}

func (s *Service) PingDB(ctx context.Context) error {
	return s.db.WithContext(ctx).Raw("SELECT 1").Error
}
