# podzol

A docker container scaling, port forwarding &amp; lifetime manager

## Usage

See `podzol help` for usage instructions.

## Configuration

Use `podzol defaultconfig` to generate a default configuration file. Edit as you see fit. Place the configuration file at `/etc/podzol/config.yaml` for the system-wide configuration.

## API Reference

All API expects JSON input and produces JSON output. Certain GET endpoints may accept query parameters.

All client commands produce their request URL and body on standard error if `-v` / `--verbose` is specified.
