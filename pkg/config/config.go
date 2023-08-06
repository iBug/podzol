package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"github.com/ustclug/podzol/pkg"
)

const (
	Format      = "yaml"
	ExampleFile = "config.example." + Format
)

func Load() error {
	err := viper.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		err = fmt.Errorf("config file not found, writing default to %s", ExampleFile)
		if err := Save(ExampleFile); err != nil {
			return err
		}
		return err
	}
	return err
}

func Save(in string) error {
	return viper.WriteConfigAs(in)
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType(Format)
	viper.SetConfigPermissions(0o644)
	viper.AddConfigPath(".") // optional
	viper.AddConfigPath(fmt.Sprintf("/etc/%s/", strings.ToLower(pkg.Name)))
	viper.SetEnvPrefix(strings.ToUpper(pkg.Name))

	viper.SetDefault("listen-addr", "127.0.0.1:9998")
	viper.SetDefault("container-prefix", strings.ToLower(pkg.Name))
	viper.SetDefault("port-min", 10000)
	viper.SetDefault("port-max", 19999)
}
