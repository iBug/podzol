package docker

import (
	"context"
	"slices"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/ustclug/podzol/pkg"
)

// Rebuild port pool from containers
func (c *Client) ResetPorts(ctx context.Context) error {
	containers, err := c.c.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("label", pkg.ID)),
	})
	if err != nil {
		return err
	}

	ports := make([]uint16, 0, len(containers))
	for _, container := range containers {
		port := container.Ports[0].PublicPort
		ports = append(ports, port)
	}
	slices.Sort(ports)

	c.pool.Load(ports)
	return nil
}
