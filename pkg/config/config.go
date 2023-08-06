package config

import (
	"strings"

	"github.com/ustclug/podzol/pkg"
)

type Config struct {
	ContainerPrefix string
	PortMin         uint16
	PortMax         uint16
}

var Default Config

func Load() error {
	Default = Config{
		ContainerPrefix: strings.ToLower(pkg.Name),
	}
	return nil
}
