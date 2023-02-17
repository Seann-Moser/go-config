package db

import (
	"context"

	"github.com/Seann-Moser/QueryHelper/dataset"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/Seann-Moser/go-config/flags"
)

const (
	datasetNameFlag     = "dataset-name"
	datasetDropTable    = "dataset-drop-table"
	datasetCreateTable  = "dataset-create-table"
	datasetUpdateColumn = "dataset-update-column"
)

func Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("db", pflag.ExitOnError)
	fs.AddFlagSet(MySqlFlags())
	fs.String(datasetNameFlag, "default", "")
	fs.Bool(datasetDropTable, false, "")
	fs.Bool(datasetCreateTable, true, "")
	fs.Bool(datasetUpdateColumn, false, "")
	return fs
}

type DB struct {
	DB            *sqlx.DB
	DatasetManger *dataset.Dataset
}

func New(ctx context.Context, logger *zap.Logger, tables ...interface{}) (*DB, error) {
	newDB := &DB{}
	db, err := connectToDB(ctx, logger)
	if err != nil {
		return nil, err
	}
	newDB.DB = db
	datasetName, err := flags.RequiredString(datasetNameFlag)
	if err != nil {
		return nil, err
	}
	ds, err := dataset.New(context.Background(), datasetName, viper.GetBool(datasetCreateTable), viper.GetBool(datasetDropTable), viper.GetBool(datasetUpdateColumn), logger, db, tables...)
	if err != nil {
		return nil, err
	}
	newDB.DatasetManger = ds
	return newDB, nil
}
