package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/go-connections/nat"
	"github.com/ustclug/podzol/pkg"
	"golang.org/x/exp/slices"
)

const Namespace = "podzol"

// ContainerOptions is the options for Create, Remove and List.
type ContainerOptions struct {
	User     int           `json:"user"`
	Token    string        `json:"token"`
	AppName  string        `json:"app"`
	Image    string        `json:"image"`
	Port     uint16        `json:"port"`
	Lifetime time.Duration `json:"lifetime"`
}

// Auxiliary struct for JSON.
type containerOptionsA ContainerOptions

// Auxiliary struct for JSON.
type containerOptionsS struct {
	containerOptionsA

	Lifetime any `json:"lifetime"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *ContainerOptions) UnmarshalJSON(b []byte) (err error) {
	aux := &containerOptionsS{containerOptionsA: containerOptionsA(*c)}
	if err = json.Unmarshal(b, aux); err != nil {
		return
	}
	switch lifetime := aux.Lifetime.(type) {
	case string:
		c.Lifetime, err = time.ParseDuration(lifetime)
	case float64:
		c.Lifetime = time.Duration(lifetime * float64(time.Second))
	default:
		err = fmt.Errorf("invalid lifetime type: %T", lifetime)
	}
	return
}

// ContainerLabel is the label data for containers.
type ContainerLabel struct {
	User     int           `json:"user"`
	App      string        `json:"challenge"`
	Lifetime time.Duration `json:"lifetime"`
}

// Auxiliary struct for JSON.
type containerLabelA ContainerLabel

// Auxiliary struct for JSON.
type containerLabelS struct {
	containerLabelA

	Lifetime string `json:"lifetime"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *ContainerLabel) UnmarshalJSON(b []byte) (err error) {
	aux := &containerLabelS{containerLabelA: containerLabelA(*c)}
	if err = json.Unmarshal(b, aux); err != nil {
		return
	}
	c.Lifetime, err = time.ParseDuration(aux.Lifetime)
	return
}

// ContainerInfo carrys the information of a container.
type ContainerInfo struct {
	Name     string    `json:"name"`
	ID       string    `json:"id"`
	Port     uint16    `json:"port"`
	Deadline time.Time `json:"deadline"`
}

// Auxiliary struct for JSON.
type containerInfoA ContainerInfo

// Auxiliary struct for JSON.
type containerInfoS struct {
	containerInfoA

	Deadline int64 `json:"deadline"`
}

// MarshalJSON implements json.Marshaler. Note that Deadline is exported as a Unix timestamp.
func (c ContainerInfo) MarshalJSON() ([]byte, error) {
	aux := &containerInfoS{containerInfoA: containerInfoA(c)}
	aux.Deadline = c.Deadline.Unix()
	return json.Marshal(aux)
}

// Construct container name from options.
func (opts *ContainerOptions) ContainerName() string {
	return fmt.Sprintf("%s_%d_%s_1", Namespace, opts.User, opts.AppName)
}

// Construct JSON data from options.
func (opts *ContainerOptions) Label() (string, error) {
	b, err := json.Marshal(ContainerLabel{
		User:     opts.User,
		App:      opts.AppName,
		Lifetime: opts.Lifetime,
	})
	return string(b), err
}

// Create a container from the given options.
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

// Remove a container.
func (c *Client) Remove(ctx context.Context, opts ContainerOptions) error {
	return c.c.ContainerRemove(ctx, opts.ContainerName(), types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   true,
		Force:         true,
	})
}

// List containers.
// Options are used to filter containers.
// Only UserID, AppName and Port are used.
func (c *Client) List(ctx context.Context, opts ContainerOptions) ([]ContainerInfo, error) {
	containers, err := c.c.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("label", pkg.ID)),
	})
	if err != nil {
		return nil, err
	}

	var infos []ContainerInfo

	for _, container := range containers {
		labelStr := container.Labels[pkg.ID]
		var label ContainerLabel
		if err := json.Unmarshal([]byte(labelStr), &label); err != nil {
			// Log error
			continue
		}

		if opts.User != 0 && label.User != opts.User {
			continue
		}
		if opts.AppName != "" && label.App != opts.AppName {
			continue
		}
		if opts.Port != 0 && !slices.ContainsFunc(container.Ports,
			func(p types.Port) bool { return p.PublicPort == opts.Port }) {
			continue
		}

		infos = append(infos, ContainerInfo{
			Name:     strings.TrimPrefix(container.Names[0], "/"),
			ID:       container.ID,
			Port:     opts.Port,
			Deadline: time.Unix(container.Created, 0).Add(label.Lifetime),
		})
	}
	return infos, nil
}
