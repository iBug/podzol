package config

import (
	"strings"

	"github.com/ustclug/podzol/pkg"
)

type Config struct {
	ContainerPrefix string
}

var Default Config

func Load() error {
	Default = Config{
		ContainerPrefix: strings.ToLower(pkg.Name),
	}
	return nil
}
