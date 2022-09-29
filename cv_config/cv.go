package cv_config

import (
	"context"
	"strings"

	"github.com/Seann-Moser/go-config/options"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func New(cmd *cobra.Command, persistent bool, options ...options.ConfigOptions) {
	for _, o := range options {
		if persistent {
			cmd.PersistentFlags().AddFlagSet(o.Flags())
		} else {
			cmd.Flags().AddFlagSet(o.Flags())
		}
	}
}

func Init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

func Execute(cmd *cobra.Command, options ...options.ConfigOptions) error {
	Init()
	New(cmd, true, options...)
	return cmd.ExecuteContext(context.Background())
}

func PersistentPreRunE(cmd *cobra.Command, _ []string) error {
	return viper.BindPFlags(cmd.Flags())
}
