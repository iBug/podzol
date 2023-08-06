package docker

import (
	"github.com/docker/docker/client"
	"github.com/ustclug/podzol/pkg/portpool"
)

type Client struct {
	c    *client.Client
	pool *portpool.Pool
}

func NewClient(portMin, portMax int) (*Client, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	return &Client{
		c:    cli,
		pool: portpool.NewPool(portMin, portMax),
	}, err
}
