package docker

import (
	"context"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/viper"
)

type Client struct {
	c      *client.Client
	prefix string

	hostnameMap     map[string]string
	hostnameMapLock sync.RWMutex
}

func NewClient(v *viper.Viper) (*Client, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	return &Client{
		c:           cli,
		prefix:      v.GetString("container-prefix"),
		hostnameMap: make(map[string]string),
	}, err
}

func (c *Client) Info(ctx context.Context) (types.Info, error) {
	return c.c.Info(ctx)
}

func (c *Client) LookupHostname(ctx context.Context) (string, bool) {
	c.hostnameMapLock.RLock()
	defer c.hostnameMapLock.RUnlock()
	ip, ok := c.hostnameMap[c.prefix]
	return ip, ok
}

func (c *Client) AddHostname(ctx context.Context, id string) error {
	ip, err := c.GetIP(ctx, id)
	if err != nil {
		return err
	}

	c.hostnameMapLock.Lock()
	defer c.hostnameMapLock.Unlock()
	c.hostnameMap[c.prefix] = ip
	return nil
}

func (c *Client) RemoveHostname(ctx context.Context, hostname string) error {
	c.hostnameMapLock.Lock()
	defer c.hostnameMapLock.Unlock()
	delete(c.hostnameMap, hostname)
	return nil
}
