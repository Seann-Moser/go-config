package flags

import (
	"time"

	"github.com/spf13/viper"
)

// InvalidFlagError is an error type used to return invalid or missing values for flags
type InvalidFlagError struct {
	Flag string
}

func (i InvalidFlagError) Error() string {
	return "flags: missing or invalid value for " + i.Flag
}

// RequiredStringSlice returns the string retrieved by viper and an error if the string is empty
func RequiredStringSlice(flag string) (s []string, err error) {
	s = viper.GetStringSlice(flag)
	if len(s) == 0 {
		err = InvalidFlagError{Flag: flag}
	}
	return
}

// RequiredString returns the string retrieved by viper and an error if the string is empty
func RequiredString(flag string) (s string, err error) {
	s = viper.GetString(flag)
	if s == "" {
		err = InvalidFlagError{Flag: flag}
	}
	return
}

// RequiredInt returns the int retrieved by viper and an error if the value is zero
func RequiredInt(flag string) (i int, err error) {
	i = viper.GetInt(flag)
	if i == 0 {
		err = InvalidFlagError{Flag: flag}
	}
	return
}

// RequiredInt64 is like RequiredInt but returns an int64
func RequiredInt64(flag string) (i int64, err error) {
	i = viper.GetInt64(flag)
	if i == 0 {
		err = InvalidFlagError{Flag: flag}
	}
	return
}

// RequiredDuration returns the duration retrieved by viper and an error if the value is zero
func RequiredDuration(flag string) (i time.Duration, err error) {
	i = viper.GetDuration(flag)
	if i == 0 {
		err = InvalidFlagError{Flag: flag}
	}
	return
}

// RequiredTime returns the time retrieved by viper and an error if the value is zero
func RequiredTime(flag string) (i time.Time, err error) {
	i = viper.GetTime(flag)
	if i.Unix() == 0 {
		err = InvalidFlagError{Flag: flag}
	}
	return
}
