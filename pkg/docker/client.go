package docker

import (
	"github.com/docker/docker/client"
	"github.com/spf13/viper"
	"github.com/ustclug/podzol/pkg/portpool"
)

type Client struct {
	c      *client.Client
	pool   *portpool.Pool
	prefix string
}

func NewClient(v *viper.Viper) (*Client, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	return &Client{
		c:      cli,
		pool:   portpool.NewPool(v.GetUint16("port-min"), v.GetUint16("port-max")),
		prefix: v.GetString("container-prefix"),
	}, err
}
