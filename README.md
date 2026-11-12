# Workplace

**A lightweight development environment manager for projects.**

Workplace is a command-line tool written in Go that helps developers configure, start, stop, inspect and manage local development environments.

It supports multiple execution backends, allowing developers to run project services using Docker, Apple Container or native host processes.

## Features

- **Multiple runtimes** — Docker, Apple Container and native execution.
- **Magento 2 support** — Preconfigured PHP-FPM, Nginx, MySQL and OpenSearch.
- **Symfony support** — Preconfigured PHP-FPM, Nginx and MySQL.
- **Configurable services** — Customize images, environment variables, ports, volumes and dependencies.
- **Local overrides** — Keep machine-specific settings separate from shared project configuration.
- **Native services** — Run locally installed services without containers.
- **Self-update** — Update Workplace directly from GitHub releases.

## Requirements

### Runtime requirements

| Runtime | Requirements |
| --- | --- |
| Docker | Docker CLI, Docker daemon and Docker Compose plugin |
| Apple Container | macOS with Apple Container CLI installed and initialized |
| Native | Required service executables installed and available in `PATH` |

The native runtime uses locally installed services rather than downloading container images.

## Installation

Prebuilt binaries are available on the GitHub Releases page.

### Linux / macOS

Download the latest release for your platform and architecture:

bash

```bash
# Linux AMD64
curl -L https://github.com/mygento/workplace/releases/latest/download/workplace_linux_amd64.tar.gz -o workplace.tar.gz

# Extract
tar -xzf workplace.tar.gz

# Install
sudo install -m 755 workplace /usr/local/bin/workplace
```

For macOS Apple Silicon, use `darwin_arm64`; for Intel Macs, use `darwin_amd64`.

### Verify installation

bash

```bash
workplace version
```

### Update

Update to the latest version:

bash

```bash
workplace self-update
```

## Quick start

Navigate to an existing Magento 2 or Symfony project containing its application files.

Create a `.workplace.yaml` file in the project root.

For example, a Magento 2 project:

yaml

```yaml
type: magento2
```

Start the environment:

bash

```bash
workplace start
```

Workplace loads the configuration, applies default services, generates the runtime configuration and starts the environment.

Stop the environment:

bash

```bash
workplace stop
```

Remove managed runtime resources:

bash

```bash
workplace delete
```

## CLI commands

### `workplace start`

Start the project's managed services.

bash

```bash
workplace start
```

Select a runtime:

bash

```bash
workplace start --runtime docker
workplace start --runtime apple-container
workplace start --runtime native
```

Specify a project directory:

bash

```bash
workplace start --dir /path/to/project
```

Docker is the default runtime.

### `workplace stop`

Stop services without deleting their resources.

bash

```bash
workplace stop
```

With a specific runtime:

bash

```bash
workplace stop --runtime native
```

### `workplace delete`

Delete the project's managed runtime resources.

bash

```bash
workplace delete
```

**Note:** Deleting resources is different from stopping services. Review the runtime's resource and volume behavior before using this command.

### `workplace inspect`

Inspect the local development environment.

bash

```bash
workplace inspect
```

Inspect another project:

bash

```bash
workplace inspect --dir /path/to/project
```

Output machine-readable JSON:

bash

```bash
workplace inspect --json
```

The inspection collects information about:

- Operating system and architecture
- Shell environment
- Installed PHP versions and PHP configuration
- Composer installation and version
- Node.js and npm versions
- Docker CLI installation
- Docker daemon availability
- Docker containers
- Project manifests and project information
- Environment findings and warnings

JSON output is useful for automation and integration with other developer tools.

### `workplace migrate-config`

Migrate legacy Workplace configuration from `package.json` to `.workplace.yaml`.

bash

```bash
workplace migrate-config
```

Specify a project directory:

bash

```bash
workplace migrate-config --dir /path/to/project
```

Overwrite an existing configuration:

bash

```bash
workplace migrate-config --overwrite
```

Without `--overwrite`, an existing `.workplace.yaml` is protected from accidental replacement.

### `workplace version`

Display the Workplace version.

bash

```bash
workplace version
```

### `workplace self-update`

Check GitHub releases and update the current executable.

bash

```bash
workplace self-update
```

The `update` alias is also available:

bash

```bash
workplace update
```

### Global verbose option

Display commands executed by Workplace:

bash

```bash
workplace --verbose start
```

The short form is:

bash

```bash
workplace -v start
```

## Configuration

Workplace uses YAML configuration instead of requiring developers to maintain all environment settings in application manifests.

There are two configuration files.

| File | Purpose |
| --- | --- |
| `.workplace.yaml` | Shared project configuration |
| `.workplace/config.local.yaml` | Local machine-specific configuration |

The shared configuration can be committed to Git. The local configuration is intended for individual developer overrides.

### Basic configuration

Magento 2:

yaml

```yaml
type: magento2
```

Symfony:

yaml

```yaml
type: symfony
```

Workplace automatically supplies the default services for the selected project type.

### Service configuration

A more complete example:

yaml

```yaml
name: my-project
type: magento2

services:
  php:
    image: php:8.3-fpm

  nginx:
    image: nginx:alpine
    port: 8081

  mysql:
    image: mysql:8.4

  redis:
    image: redis:alpine

  rabbitmq:
    image: rabbitmq:management
```

### Configuration reference

The top-level configuration supports:

| Property | Description |
| --- | --- |
| `name` | Project name used for runtime resource naming |
| `type` | Project type: `magento2` or `symfony` |
| `services` | Service definitions and overrides |

Each service supports:

| Property | Description |
| --- | --- |
| `image` | Container image |
| `port` | Host port |
| `env` | Environment variables |
| `volumes` | Volume mount definitions |
| `depends_on` | Service dependencies |
| `container_name` | Container name |
| `managed` | Whether the native runtime manages the service (local configuration only) |

Volume configuration:

yaml

```yaml
services:
  mysql:
    volumes:
      - name: my-project-db
        target: /var/lib/mysql
        type: volume
```

Bind mount example:

yaml

```yaml
services:
  php:
    volumes:
      - source: .
        target: /var/www/magento2
        type: bind
```

Read-only bind mount:

yaml

```yaml
services:
  nginx:
    volumes:
      - source: .
        target: /var/www/magento2
        type: bind
        read_only: true
```

### Default services

Workplace provides the following defaults:

| Service | Default image | Magento 2 | Symfony |
| --- | --- | --- | --- |
| PHP | `php:8.3-fpm` | Yes | Yes |
| Nginx | `nginx:alpine` | Yes | Yes |
| MySQL | `mysql:8.4` | Yes | Yes |
| OpenSearch | `opensearchproject/opensearch:3` | Yes | Optional |
| Redis | `redis:alpine` | Optional | Optional |
| Varnish | `varnish:stable` | Optional | Optional |
| RabbitMQ | `rabbitmq:management` | Optional | Optional |

Optional services are added when explicitly declared in the configuration.

### Local configuration

Machine-specific configuration should be stored in:

text

```text
.workplace/
└── config.local.yaml
```

For example:

yaml

```yaml
services:
  php:
    image: php:8.3-fpm

  redis:
    managed: false
```

The `managed` property is available only in local configuration.

Setting `managed: false` tells the native runtime not to start, stop or delete that service.

Local configuration cannot override the project name or project type.