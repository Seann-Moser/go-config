package logging

import "github.com/spf13/pflag"

const (
	productionFlag   = "production"
	loggingLevelFlag = "logging-level"
	nameFlag         = "logger-name"
)

func Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("logging", pflag.ExitOnError)
	fs.Bool(productionFlag, false, "sets zap production data")
	fs.String(loggingLevelFlag, "info", "")
	fs.String(nameFlag, "default", "")
	return fs
}
