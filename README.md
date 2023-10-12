# podzol

A docker container scaling, port forwarding &amp; lifetime manager

## Usage

See `podzol help` for usage instructions.

### Configuration

Use `podzol defaultconfig` to generate a default configuration file. Edit as you see fit. Place the configuration file at `/etc/podzol/config.yaml` for the system-wide configuration.

### Deployment

Please run the server using `127.0.0.1:port` as listen address and place Nginx or Apache2 in front of it. Then you can configure SSL/TLS and access control with Nginx.

## API Reference

All API expects JSON input and produces JSON output. It is always recommended to set `Content-Type: application/json`. Certain GET endpoints may accept query parameters.

All client commands produce their request URL and body on standard error if `-v` / `--verbose` is specified.

### Base types

Base request type:

```go
type ContainerOptions struct {
    // User ID
	User     int           `json:"user"`

    // Token to be supplied to the container
	Token    string        `json:"token"`

    // For identification purposes
	AppName  string        `json:"app"`

    // First segment of the Host header, for reverse proxying
	Hostname string        `json:"hostname"`

    // Docker image to be used
	Image    string        `json:"image"`

    // How long should podzol auto-destroy the container, in seconds
	Lifetime time.Duration `json:"lifetime"`
}
```

Base response type:

```go
type ContainerInfo struct {
    // Container name, ID, reverse proxy hostname
	Name     string    `json:"name"`
	ID       string    `json:"id"`
	Hostname string    `json:"hostname"`

    // When the container will expire, in Unix timestamp
	Deadline time.Time `json:"deadline"`
}
```

### Create container

```
POST /create
```

All fields are required.

Returns a single `ContainerInfo` struct.

### Remove container

```
POST /remove
```

Only `user` and `app` fields are required.

Returns an empty object.

### List containers

```
GET /list?opts=...
POST /list
```

`opts` is a JSON-encoded `ContainerOptions` struct. Only `user` and `app` fields are respected, if supplied.

Returns a list of `ContainerInfo` structs.

### Purge containers

This endpoint purges all "expired" containers.

```
POST /purge
```

No body is required.

Returns a list of `ContainerInfo` structs for the containers that have been attempted to remove.

Usually this endpoint is not called by an application, but rather by a cron job.

## Known Issues

- None
