# Performance & Overhead

`dir-cpu` is designed to be a lightweight observer — it should not meaningfully affect the system it's measuring.

## What happens each cycle

1. **`process.Processes()`** — enumerates all PIDs from `/proc`. One `readdir` on a virtual filesystem.
2. **`p.Times()`** per process — reads `/proc/[pid]/stat` for each PID. One file read per process.
3. **`p.Cwd()` or `p.Exe()`** per process — resolves a symlink in `/proc/[pid]/`. One `readlink` per process.
4. **`time.Sleep(interval)`** — the dominant cost by far.
5. **Delta computation + path roll-up** — pure in-memory arithmetic.
6. **Terminal render** — one write to stdout.

Steps 1–3 are I/O against `/proc`, which is an in-memory virtual filesystem. Reads do not touch physical disk.

## Syscall budget

On a system with 300 processes:

| Operation | Syscalls |
|-----------|----------|
| `readdir /proc` | 1 |
| `read /proc/[pid]/stat` × 300 | ~300 |
| `readlink /proc/[pid]/cwd` × 300 | ~300 |
| **Total per cycle** | **~601** |

At a 1-second interval, this is ~601 syscalls per second. For comparison, a busy web server handles thousands of syscalls per millisecond.

## CPU usage of dir-cpu itself

On a quiet desktop system (Ubuntu 22.04, ~250 processes, 1s interval):

- `dir-cpu` itself: **< 1% CPU**
- Memory: **~10–15 MB RSS** (dominated by the gopsutil process list and the map allocations)

On a very busy server (1000+ processes):

- CPU rises roughly linearly with process count
- At 1000 processes: still typically **< 3% CPU** on a modern machine

## GC pressure

The main allocation per cycle is:

- A `map[int32]procData` with one entry per process (new each cycle)
- A `map[string]float64` for directory aggregation
- A `[]entry` slice for sorting

These are all short-lived and reclaimed by Go's GC between cycles. The sleep period gives the GC ample time to run without pausing the display.

## Reducing overhead

If overhead is a concern on a very busy system:

```bash
# Longer interval = fewer /proc reads per minute
dir-cpu -i 5s

# Higher threshold = fewer path traversals (short-circuit on low-CPU processes)
dir-cpu -t 2.0
```

The `-t` flag is particularly effective: if 80% of processes use < 0.5% CPU, they're skipped entirely before any path traversal happens.

## What dir-cpu does NOT do

- No persistent background daemon
- No kernel module or eBPF hooks
- No disk I/O (all reads are against `/proc`, in-memory)
- No network activity
- No modification of system state
