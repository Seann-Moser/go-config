package logging

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func New() (*zap.Logger, error) {
	logger, err := setupLogger(viper.GetBool(productionFlag), viper.GetString(loggingLevelFlag))
	if err != nil {
		return nil, err
	}
	logger.Named(viper.GetString(nameFlag))
	return logger, nil
}

func setupLogger(production bool, level string) (*zap.Logger, error) {
	var conf zap.Config
	if production {
		conf = zap.NewProductionConfig()
	} else {
		conf = zap.NewDevelopmentConfig()
	}

	if err := conf.Level.UnmarshalText([]byte(level)); err != nil {
		return nil, err
	}

	return conf.Build()
}
