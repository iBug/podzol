# podzol

A docker container scaling, port forwarding &amp; lifetime manager

## Usage

See `podzol help` for usage instructions.

### Configuration

Use `podzol defaultconfig` to generate a default configuration file. Edit as you see fit. Place the configuration file at `/etc/podzol/config.yaml` for the system-wide configuration.

### Deployment

Please run the server using `127.0.0.1:port` as listen address and place Nginx or Apache2 in front of it. Then you can configure SSL/TLS and access control with Nginx.

## API Reference

All API expects JSON input and produces JSON output. Certain GET endpoints may accept query parameters.

All client commands produce their request URL and body on standard error if `-v` / `--verbose` is specified.

### Known Issues

- None known.
