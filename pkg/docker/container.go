package docker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/ustclug/podzol/pkg"
)

// ContainerOptions is the options for Create, Remove and List.
type ContainerOptions struct {
	User     int           `json:"user"`
	Token    string        `json:"token"`
	AppName  string        `json:"app"`
	Hostname string        `json:"hostname"`
	Image    string        `json:"image"`
	Lifetime time.Duration `json:"lifetime"`
}

// Auxiliary struct for JSON.
type containerOptionsA ContainerOptions

// Auxiliary struct for JSON.
type containerOptionsS struct {
	*containerOptionsA

	Lifetime any `json:"lifetime"`
}

// MarshalJSON implements json.Marshaler.
func (c ContainerOptions) MarshalJSON() ([]byte, error) {
	aux := &containerOptionsS{containerOptionsA: (*containerOptionsA)(&c)}
	aux.Lifetime = int64(c.Lifetime / time.Second)
	return json.Marshal(aux)
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *ContainerOptions) UnmarshalJSON(b []byte) (err error) {
	aux := &containerOptionsS{containerOptionsA: (*containerOptionsA)(c)}
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
	*containerLabelA

	Lifetime string `json:"lifetime"`
}

// MarshalJSON implements json.Marshaler.
func (c ContainerLabel) MarshalJSON() ([]byte, error) {
	aux := &containerLabelS{containerLabelA: (*containerLabelA)(&c)}
	aux.Lifetime = c.Lifetime.String()
	return json.Marshal(aux)
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *ContainerLabel) UnmarshalJSON(b []byte) (err error) {
	aux := &containerLabelS{containerLabelA: (*containerLabelA)(c)}
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
	Hostname string    `json:"hostname"`
	Deadline time.Time `json:"deadline"`
}

// Auxiliary struct for JSON.
type containerInfoA ContainerInfo

// Auxiliary struct for JSON.
type containerInfoS struct {
	*containerInfoA

	Deadline int64 `json:"deadline"`
}

// MarshalJSON implements json.Marshaler. Note that Deadline is exported as a Unix timestamp.
func (c ContainerInfo) MarshalJSON() ([]byte, error) {
	aux := &containerInfoS{containerInfoA: (*containerInfoA)(&c)}
	aux.Deadline = c.Deadline.Unix()
	return json.Marshal(aux)
}

// UnmarshalJSON implements json.Unmarshaler. Note that Deadline is expected as a Unix timestamp.
func (c *ContainerInfo) UnmarshalJSON(data []byte) error {
	aux := &containerInfoS{containerInfoA: (*containerInfoA)(c)}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	c.Deadline = time.Unix(aux.Deadline, 0)
	return nil
}

// Construct container name from options.
func (c *Client) ContainerName(opts ContainerOptions) string {
	return fmt.Sprintf("%s_%d_%s_1", c.prefix, opts.User, opts.AppName)
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
func (c *Client) Create(ctx context.Context, opts ContainerOptions) (ContainerInfo, error) {
	label, err := opts.Label()
	if err != nil {
		return ContainerInfo{}, err
	}

	containerName := c.ContainerName(opts)

	containerConfig := &container.Config{
		Hostname: containerName,
		Image:    opts.Image,
		Labels:   map[string]string{pkg.ID: label},
	}

	hostConfig := &container.HostConfig{
		NetworkMode: "bridge",
		AutoRemove:  true,
	}

	createTime := time.Now().Truncate(time.Second)
	resp, err := c.c.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, containerName)
	if err != nil {
		return ContainerInfo{}, err
	}
	if err := c.c.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		// Remove container if start failed
		_ = c.Remove(ctx, opts)
		return ContainerInfo{}, err
	}
	c.AddHostname(ctx, resp.ID)

	return ContainerInfo{
		Name:     containerName,
		ID:       resp.ID,
		Hostname: opts.Hostname,
		Deadline: createTime.Add(opts.Lifetime),
	}, err
}

func (c *Client) remove(ctx context.Context, name string) error {
	return c.c.ContainerRemove(ctx, name, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
}

// Remove a container.
func (c *Client) Remove(ctx context.Context, opts ContainerOptions) error {
	return c.remove(ctx, c.ContainerName(opts))
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

	infos := make([]ContainerInfo, 0)

	for _, container := range containers {
		labelStr := container.Labels[pkg.ID]
		var label ContainerLabel
		if err := json.Unmarshal([]byte(labelStr), &label); err != nil {
			// Log error
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		if opts.User != 0 && label.User != opts.User {
			continue
		}
		if opts.AppName != "" && label.App != opts.AppName {
			continue
		}

		infos = append(infos, ContainerInfo{
			Name:     strings.TrimPrefix(container.Names[0], "/"),
			ID:       container.ID,
			Deadline: time.Unix(container.Created, 0).Add(label.Lifetime),
		})
	}
	return infos, nil
}

// Get container IP address.
func (c *Client) GetIP(ctx context.Context, name string) (string, error) {
	inspect, err := c.c.ContainerInspect(ctx, name)
	if err != nil {
		return "", err
	}
	return inspect.NetworkSettings.IPAddress, nil
}

type ContainerActionError struct {
	Action    string        `json:"action"`
	Container ContainerInfo `json:"container"`
	Err       error         `json:"error"`
}

func (e ContainerActionError) Error() string {
	return fmt.Sprintf("%s %s: %v", e.Action, e.Container.Name, e.Err)
}

// Purge expired containers.
// Returns the list of (attempted) purged containers.
// Note that if the metadata of a container is corrupted, it will be removed as well.
// The returned error is a list of errors that occurred during the purge.
func (c *Client) Purge(ctx context.Context) ([]ContainerInfo, error) {
	containers, err := c.c.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("label", pkg.ID)),
	})
	if err != nil {
		return nil, err
	}

	infos := make([]ContainerInfo, 0)

	for _, container := range containers {
		labelStr := container.Labels[pkg.ID]
		var label ContainerLabel
		if err := json.Unmarshal([]byte(labelStr), &label); err != nil {
			// TODO: Log error
			label.Lifetime = 0
		}

		info := ContainerInfo{
			Name:     strings.TrimPrefix(container.Names[0], "/"),
			ID:       container.ID,
			Deadline: time.Unix(container.Created, 0).Add(label.Lifetime),
		}
		if time.Now().After(info.Deadline) {
			infos = append(infos, info)
		}
	}

	errs := make([]error, 0)
	for _, container := range infos {
		err := c.remove(ctx, container.Name)
		if err != nil {
			errs = append(errs, ContainerActionError{
				Action:    "remove",
				Container: container,
				Err:       err,
			})
		}
	}
	return infos, errors.Join(errs...)
}
