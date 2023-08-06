package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/ustclug/podzol/pkg"
)

const Namespace = "podzol"

type ContainerOptions struct {
	UserID        int
	Token         string
	ChallengeName string
	Image         string
	Port          int
	Lifetime      time.Duration
}

func (opts *ContainerOptions) ContainerName() string {
	return fmt.Sprintf("%s_%d_%s_1", Namespace, opts.UserID, opts.ChallengeName)
}

func (opts *ContainerOptions) Label() (string, error) {
	b, err := json.Marshal(map[string]any{
		"user":      opts.UserID,
		"challenge": opts.ChallengeName,
		"lifetime":  opts.Lifetime.String(),
	})
	return string(b), err
}

func (c *Client) Create(ctx context.Context, opts ContainerOptions) (string, error) {
	label, err := opts.Label()
	if err != nil {
		return "", err
	}

	containerConfig := &container.Config{
		Hostname: opts.ContainerName(),
		ExposedPorts: nat.PortSet{
			nat.Port(fmt.Sprintf("%d/tcp", opts.Port)): {},
		},
		Image:  opts.Image,
		Labels: map[string]string{pkg.ID: label},
	}

	hostConfig := &container.HostConfig{
		NetworkMode: "bridge",
		AutoRemove:  true,
	}

	resp, err := c.c.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, opts.ContainerName())
	if err != nil {
		return "", err
	}
	if err := c.c.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		// Remove container if start failed
		_ = c.Remove(ctx, opts)
		return "", err
	}
	return resp.ID, err
}

func (c *Client) Remove(ctx context.Context, opts ContainerOptions) error {
	return c.c.ContainerRemove(ctx, opts.ContainerName(), types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   true,
		Force:         true,
	})
}
