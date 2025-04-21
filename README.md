# Workmake Task API

A scalable service for executing heavy I/O bound background tasks written in Go.

## Table of Contents

- [Project Structure](#project-structure)
- [Makefile](#makefile)
- [Running the Application](#running-the-application)
- [API Documentation](#api-documentation)
- [Application Configuration](#application-configuration)
- [License](#license)

## Project Structure

Here are the main components of the application:

```bash
.
├── api                # API documentation
├── cmd                # Application entrypoint
└── internal
    ├── api            # Data delivery layer
    ├── config         # Configuration loading logic
    ├── entity         # Core domain entities
    ├── executor       # Tasks execution logic
    ├── handler        # Handling tasks logic
    └── manager        # Managing tasks logic
```

## Makefile

Explore avaliable `Makefile` targets:

```bash
make help
```

## Running the Application

1. Prepare config file.

2. Run the application:

    ```bash
    # You can ovveride CONFIG_PATH (Default: .config.yml)
    make run CONFIG_PATH=<path>
    ```

## API Documentation

The application is documented with Swagger. You can explore API in `api/swagger.yml`.

## Application Configuration

The application is configured via YAML file. Application uses `-configPath` flag to load configuration from YAML file (Default: `.config.yml`).

```bash
# You can override default value
make run CONFIG_PATH=<path>
```

Here is the basic structure of the configuration file:

```yaml
env: dev

task_cleanup: 30m
queue_size: 1000
max_workers: 100

server:
  port: 8080
  read_timeout: 5s
  write_timeout: 10s
  idle_timeout: 1m
  max_header_bytes: 1048576 # 1 << 20
  shutdown_timeout: 10s
```

The behavior of the application depends on the `env` passed in the configuration file:

1. `dev` - logging is structured with plain text (debug level).
2. `stage` - logging is structured with JSON (debug level).
3. `prod` - logging is structured with JSON (info level).

## License

This project is licensed under the WTFPL License - see the `LICENSE` file for details.
