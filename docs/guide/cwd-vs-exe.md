# cwd vs exe mode

`dir-cpu` can attribute CPU usage to a process in two ways. Choosing the right one depends on what you're trying to monitor.

## The fundamental problem

The OS identifies processes by PID, not by project. A single `python3` process running your web app and another `python3` process running an unrelated cron job look identical by name. `dir-cpu` solves this by using the filesystem path instead.

But there are two paths that could represent a process:

- **Where it's running from** — the working directory (`cwd`)
- **What binary is running** — the executable path (`exe`)

These are different, and each has a different sweet spot.

## cwd mode (default)

**Path used:** `/proc/[pid]/cwd`

This is the directory the process was launched from — typically the project root when you run `python app.py` or `node server.js` from a terminal.

```
/home/user/projects/myapp$ python app.py
```

In this case `cwd` = `/home/user/projects/myapp`, which is exactly where you want the cost attributed.

**Best for:**
- Python scripts
- Node.js apps
- Ruby, PHP, shell scripts
- Any interpreted language where the binary (`/usr/bin/python3`) is shared but the project directory distinguishes them

**Limitation:**  
If a process changes its working directory after startup (e.g., a daemon that does `chdir("/")` on launch), the initial project context is lost. System daemons often do this.

## exe mode

**Path used:** `/proc/[pid]/exe`

This is the path to the actual binary being executed — resolved through the symlink in `/proc`.

```
/usr/local/bin/myapp
/home/user/go/bin/server
/opt/app/bin/worker
```

**Best for:**
- Compiled Go, Rust, C, C++ programs
- Binaries installed to a specific project or environment path
- Identifying processes by where their binary lives, not where they were launched

**Limitation:**  
For interpreted languages (`/usr/bin/python3`, `/usr/bin/node`), all Python or Node processes roll up to `/usr/bin` — completely losing project context. Use `cwd` instead.

## Side-by-side comparison

| Scenario | Recommended mode |
|----------|-----------------|
| Python/Node/Ruby project | `cwd` |
| Go/Rust/C binary | `exe` |
| Mix of both | Run two terminals with each mode |
| Daemon monitoring | `exe` (daemons often `chdir /`) |
| Security scan for rogue processes | `exe` (surfaces unusual binary locations) |
| Monorepo with multiple services | `cwd` |

## Example: monorepo

```
/home/user/monorepo/
  frontend/    ← npm dev server
  backend/     ← go run .
  worker/      ← python worker.py
```

With `cwd` mode, each service costs roll up correctly to their subdirectory and then aggregate at `/home/user/monorepo`.

With `exe` mode, the frontend rolls to `/home/user/.nvm/...`, the backend to `/home/user/go/bin/...`, and the worker to `/usr/bin/` — the monorepo structure is invisible.

## Mixing both modes

There's no combined mode, but you can run two instances in split panes:

```bash
# pane 1
dir-cpu -base=cwd

# pane 2
dir-cpu -base=exe
```
