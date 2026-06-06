package log

import (
	"github.com/mkfolder/moxie/internal/common"
	"go.uber.org/zap"
)

func New(env common.Environment) (*zap.SugaredLogger, error) {
	var (
		logger *zap.Logger
		err    error
	)

	if env.IsDevelopment() {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}
