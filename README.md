# raii

Start resources when you enter a directory. Stop them when you leave.

## What it does

Manages external resources via `direnv` hooks.

When you `cd` into a project, `raii` runs a start command. When you `cd` out or close your shell, it runs a stop command.

Multiple shells in the same project share resources. Resources only die when the last shell exits.

## Install

```bash
go install github.com/maxdikun/raii@latest
```

## Usage

### 1. Config

Create `raii.toml`:

```toml
session = "my-project"

[commands]
start = "docker compose up -d"
stop  = "docker compose down"
check = "docker compose ps | grep -q Up"
```

- `session`: Resource group ID. Defaults to the config file's directory.
- `start`: Runs if resources are not running and no owners exist.
- `stop`: Runs when the last owner exits.
- `check`: Must exit 0 if resources are up.

### 2. direnv hook

Copy the helper:

```bash
mkdir -p ~/.config/direnv/lib
cp direnv/use_raii.sh ~/.config/direnv/lib/
```

Add to `.envrc`:

```bash
use raii
```

Allow it:

```bash
direnv allow
```

Done.

## How it works

`raii` stores state in `~/.local/share/raii/state.json`.

- `raii start`: Adds your shell PID to the owner list. Starts resources if needed.
- `raii stop`: Removes your PID. Stops resources if the list is empty.
- `raii check`: Runs the check command.

A background watchdog is spawned on `start`. It polls your shell PID. When your shell dies, the watchdog calls `raii stop` for you. This handles terminal closes and crashes.

## CLI

```
raii <start|stop|check|watch> [flags]

Flags:
  --config string   Config path (default "raii.toml")
  --owner string    Owner ID (default: parent PID)
```

## State directory

Override with `RAII_STATE_DIR`:

```bash
export RAII_STATE_DIR=/tmp/raii
```

## Example: not Docker

```toml
session = "dev-server"

[commands]
start = "./scripts/start.sh"
stop  = "./scripts/stop.sh"
check = "curl -sf http://localhost:8080/health"
```

## License

MIT
