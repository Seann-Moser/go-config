package db

import (
	"context"

	"github.com/Seann-Moser/QueryHelper/v2/dataset"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/Seann-Moser/go-config/flags"
)

const (
	datasetNameFlag    = "dataset-name-flag"
	datasetDropTable   = "dataset-drop-table-flag"
	datasetCreateTable = "dataset-create-table-flag"
)

func Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("db", pflag.ExitOnError)
	fs.AddFlagSet(MySqlFlags())
	fs.String(datasetNameFlag, "default", "")
	fs.Bool(datasetDropTable, false, "")
	fs.Bool(datasetCreateTable, true, "")
	return fs
}

type DB struct {
	DB            *sqlx.DB
	DatasetManger *dataset.Dataset
}

func New(logger *zap.Logger, tables ...interface{}) (*DB, error) {
	newDB := &DB{}
	db, err := connectToDB(logger)
	if err != nil {
		return nil, err
	}
	newDB.DB = db
	datasetName, err := flags.RequiredString(datasetNameFlag)
	if err != nil {
		return nil, err
	}
	ds, err := dataset.NewDataset(context.Background(), datasetName, viper.GetBool(datasetCreateTable), viper.GetBool(datasetDropTable), logger, db, tables...)
	if err != nil {
		return nil, err
	}
	newDB.DatasetManger = ds
	return newDB, nil
}
