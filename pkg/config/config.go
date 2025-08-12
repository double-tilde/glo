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
	Color string
	IgnoredDirs []string
	LogMessages bool
	Shape string
}

func GetIgnoredDirs(original, added, remove []string) []string {
	var total []string

	total = append(total, original...)
	total = append(total, added...)

	for k, t := range total {
		for _, r := range remove {
			if t == r {
				total = append(total[:k], total[k+1:]...)
			}
		}
	}

	return total
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
		Color: viper.GetString("color"),
		IgnoredDirs: GetIgnoredDirs(
			viper.GetStringSlice("default_ignored_directories"),
			viper.GetStringSlice("user_added_ignored_directories"),
			viper.GetStringSlice("user_excluded_ignored_directories"),
		),
		LogMessages: viper.GetBool("log_messages"),
		Shape: viper.GetString("shape"),
	}
}
