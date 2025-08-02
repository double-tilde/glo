package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	CommandLines     = 4
	CommandSeperator = "----"
	GitCommand       = []string{
		"git",
		"log",
		"--pretty=format:%H\n%an\n%ad\n%s\n" + CommandSeperator,
		"--date=format:%Y-%m-%d %H:%M",
	}
	GloCommitsFile   = "glo.json"
	GitDirectory     = ".git"
	GloDirectory     = "glo"
	GloHomeDirectory = "." + GloDirectory
	LogFileName      = "glo.log"
	TimeFormat       = "2006-01-02 15:04"
)

type Config struct {
	IgnoreDirs      []string
	LogMessages bool
}

func Setup(configHome string) error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	// viper.AddConfigPath("$DATA_CONFIG_HOME/" + GloDirectory)
	// viper.AddConfigPath("/etc/" + GloDirectory)
	viper.AddConfigPath(configHome)
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("fatal error config file: %w", err)
	}

	return nil
}

func New() *Config {
	return &Config{
		IgnoreDirs:      viper.GetStringSlice("ignore-directories"),
		LogMessages: viper.GetBool("log-messages"),
	}
}
