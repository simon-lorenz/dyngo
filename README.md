# DynGO

A DynDNS Client written in Go.

## Features

- 🌍 Speaks IPv4 and IPv6
- ⚙️ YAML-based configuration
- 🥇 Single binary deployment
- 📗 Supports different DNS service providers
- 🔎 Multiple strategies to determine your wan ip address

> ⚠️ DynGO is neither very mature nor feature rich at the moment. I built it for myself, but maybe you'll find it useful too. If you would like to have a currently unsupported service provider or feature implemented, please open an [issue](https://github.com/simon-lorenz/dyngo/issues/new) or send a PR.

## Installation

### Download

You can download the most recent version of DynGO from the [releases page](https://github.com/simon-lorenz/dyngo/releases).

```bash
curl -L https://github.com/simon-lorenz/dyngo/releases/latest/download/dyngo-$(uname -m) -o /usr/bin/dyngo
chmod +x /usr/bin/dyngo
```

### Configuration

DynGO is configured via YAML. You must prepare your configuration file before running DynGO. For available options see [config/example.yaml](config/example.yaml).

The default configuration path is `/etc/dyngo/config.yaml`. If you need to specify another location, use `dyngo --config=<path>`.

### Deployment

DynGO is meant to run continuously. Use a tool of your choice to ensure that DynGO is running.

#### systemd

There's an example unit at [systemd/dyngo.service](systemd/dyngo.service).
