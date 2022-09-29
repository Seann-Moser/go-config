package options

import "github.com/spf13/pflag"

type ConfigOptions interface {
	Flags() *pflag.FlagSet
}
