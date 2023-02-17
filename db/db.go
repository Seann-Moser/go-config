package db

import (
	"context"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/pflag"
	"go.uber.org/zap"

	"github.com/Seann-Moser/go-config/flags"
)

const (
	sqlDbHostFlag       = "sql-db-host"
	sqlDbPortFlag       = "sql-db-port"
	sqlDbPasswordFlag   = "sql-db-password"
	sqlDbUserFlag       = "sql-db-user"
	sqlDbMaxConnections = "sql-db-max-connections"
)

func MySqlFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("sqldb", pflag.ExitOnError)
	fs.String(sqlDbHostFlag, "", "")
	fs.String(sqlDbPortFlag, "", "")
	fs.String(sqlDbPasswordFlag, "", "")
	fs.String(sqlDbUserFlag, "", "")
	fs.Int(sqlDbMaxConnections, 10, "")
	return fs
}

func connectToDB(ctx context.Context, logger *zap.Logger) (*sqlx.DB, error) {
	host, err := flags.RequiredString(sqlDbHostFlag)
	if err != nil {
		return nil, err
	}
	port, err := flags.RequiredString(sqlDbPortFlag)
	if err != nil {
		return nil, err
	}
	password, err := flags.RequiredString(sqlDbPasswordFlag)
	if err != nil {
		return nil, err
	}
	user, err := flags.RequiredString(sqlDbUserFlag)
	if err != nil {
		return nil, err
	}
	maxConnections, err := flags.RequiredInt(sqlDbMaxConnections)
	if err != nil {
		return nil, err
	}
	dbConf := mysql.Config{
		AllowNativePasswords: true,
		User:                 user,
		Passwd:               password,
		Net:                  "tcp",
		Addr:                 host + ":" + port,
		CheckConnLiveness:    true,
	}
	logger.Debug("connecting to sql db", zap.String("dsn", dbConf.FormatDSN()))

	db, err := sqlx.Open("mysql", dbConf.FormatDSN())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxConnections) //max open connections
	ticker := time.NewTicker(1 * time.Second)
	var retries int
	for {
		select {
		case <-ctx.Done():
			return nil, err
		case <-ticker.C:
			logger.Debug("attempting to connect to db", zap.Int("attempt", retries))
			if err = db.Ping(); err == nil {
				return db, nil
			}
			retries++
		}
	}
}
