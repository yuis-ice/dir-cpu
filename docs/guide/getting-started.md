# Getting Started

## Requirements

- **Linux** — reads `/proc`, so macOS and Windows are not supported
- **Go 1.21+** — only needed if building from source

## Install

### From source

```bash
git clone https://github.com/yuis-ice/dir-cpu
cd dir-cpu
go build -o dir-cpu .
sudo mv dir-cpu /usr/local/bin/
```

### Run without installing

```bash
go run github.com/yuis-ice/dir-cpu@latest
```

### Verify

```bash
dir-cpu --help
```

Expected output:

```
Usage of dir-cpu:
  -base string   aggregation basis: cwd | exe  (default "cwd")
  -i duration    update interval               (default 1s)
  -t float       display threshold (%)         (default 0.5)
  -n int         max rows to display           (default 40)
```

## First run

```bash
dir-cpu
```

After one interval (default 1 second), the screen clears and you'll see something like:

```
CPU% by cwd dir  (updated: 14:32:01)
────────────────────────────────────────
   45.2%  /home/user/projects/myapp
   45.2%  /home/user/projects
   45.2%  /home/user
   12.1%  /home/user/projects/myapp/worker
    4.7%  /usr/lib
```

Each line shows the **total CPU** of all processes whose working directory is inside that path.

## Quick examples

```bash
# Default: cwd-based, 1s interval, hide anything below 0.5%
dir-cpu

# Faster refresh
dir-cpu -i 500ms

# Show only the top 10 rows
dir-cpu -n 10

# Group by executable path instead of working directory
dir-cpu -base=exe

# Lower threshold to catch light background tasks
dir-cpu -t 0.1
```

## Next steps

- [CLI Reference](./cli-reference) — every flag explained
- [cwd vs exe mode](./cwd-vs-exe) — which mode is right for your use case
- [How it works](../architecture/overview) — the design behind the numbers
