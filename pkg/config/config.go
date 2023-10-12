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
	return viper.ReadInConfig()
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
	viper.SetDefault("http-addr", "127.0.0.1:9999")
	viper.SetDefault("container-prefix", strings.ToLower(pkg.Name))
}
