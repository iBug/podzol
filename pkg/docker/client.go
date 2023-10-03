package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/viper"
)

type Client struct {
	c      *client.Client
	prefix string
}

func NewClient(v *viper.Viper) (*Client, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	return &Client{
		c:      cli,
		prefix: v.GetString("container-prefix"),
	}, err
}

func (c *Client) Info(ctx context.Context) (types.Info, error) {
	return c.c.Info(ctx)
}
