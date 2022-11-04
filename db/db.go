package db

import (
	"time"

	"github.com/Seann-Moser/go-config/flags"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

const (
	sqlDbHostFlag       = "sql-db-host"
	sqlDbPortFlag       = "sql-db-port"
	sqlDbPasswordFlag   = "sql-db-password"
	sqlDbUserFlag       = "sql-db-user"
	sqlDbWaitToConnect  = "sql-db-wait-to-connect"
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

func connectToDB(logger *zap.Logger) (*sqlx.DB, error) {
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
	var retries int
	for {
		logger.Debug("attempting to connect to db", zap.Int("attempt", retries))
		time.Sleep(1 * time.Second)
		if err := db.Ping(); err == nil {
			break
		}
		retries++
	}
	return db, err
}
