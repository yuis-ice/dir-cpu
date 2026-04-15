# dir-cpu

[![CI](https://github.com/yuis-ice/dir-cpu/actions/workflows/ci.yml/badge.svg)](https://github.com/yuis-ice/dir-cpu/actions/workflows/ci.yml)
[![Docs](https://img.shields.io/badge/docs-yuis--ice.github.io%2Fdir--cpu-blue)](https://yuis-ice.github.io/dir-cpu/)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go Report](https://goreportcard.com/badge/github.com/yuis-ice/dir-cpu)](https://goreportcard.com/report/github.com/yuis-ice/dir-cpu)

**[Full documentation →](https://yuis-ice.github.io/dir-cpu/)**

A real-time CLI that shows CPU usage aggregated by **filesystem directory** — not by process name or cgroup.

If `/home/user/projects/myapp/server.py` uses 30% CPU, `dir-cpu` shows that cost against every ancestor:

```
   30.0%  /home/user/projects/myapp
   30.0%  /home/user/projects
   30.0%  /home/user
```

Works on Linux. Useful for spotting which project directory is burning CPU without cross-referencing process names manually.

## Install

**From source (requires Go 1.21+):**

```bash
git clone https://github.com/yuis-ice/dir-cpu
cd dir-cpu
go build -o dir-cpu .
sudo mv dir-cpu /usr/local/bin/
```

**Run without installing:**

```bash
go run github.com/yuis-ice/dir-cpu@latest
```

## Usage

```
dir-cpu [flags]

Flags:
  -base string   aggregation basis: cwd | exe  (default "cwd")
  -i duration    update interval               (default 1s)
  -t float       display threshold (%)         (default 0.5)
  -n int         max rows to display           (default 40)
```

**cwd mode** (default) — groups by each process's working directory. Best for scripts and interpreted languages (`python`, `node`, `ruby`) run from a project folder.

**exe mode** — groups by the directory containing each process's binary. Better for compiled programs.

### Examples

```bash
# Watch which project directory is consuming the most CPU
dir-cpu

# exe-based, faster refresh, lower threshold
dir-cpu -base=exe -i 500ms -t 0.1

# Show only top 10 directories
dir-cpu -n 10
```

## How it works

`dir-cpu` takes two snapshots of `/proc/[pid]/stat` CPU times separated by the configured interval, computes the delta per process, then rolls each process's usage up through every ancestor directory.

The result at any directory is the **sum of CPU% of all processes currently running inside it** — the same metric `top` shows per process, but bucketed by path hierarchy.

Percentages can exceed 100% on multi-core systems (consistent with `top -H` behavior).

## Permissions

Running as a regular user shows only your own processes. For full system visibility, run with `sudo`.

## Requirements

- Linux (reads `/proc`)
- Go 1.21+ (to build from source)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT
