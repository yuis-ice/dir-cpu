# Tree Roll-up Algorithm

The core insight of `dir-cpu` is that a process's CPU cost belongs not only to its immediate directory, but to every ancestor up to the root. This is the **roll-up**.

## How it works

For each process with a non-zero CPU delta, `dir-cpu` walks up the path hierarchy and adds the CPU percentage to each directory along the way.

```go
dir := filepath.Clean(path)   // e.g. "/home/user/projects/app"
for {
    dirCPU[dir] += cpuPct
    if dir == "/" {
        break
    }
    parent := filepath.Dir(dir)
    if parent == dir {         // reached filesystem root
        break
    }
    dir = parent
}
```

### Example

Process at `/home/user/projects/app/worker` using 40% CPU:

```
/home/user/projects/app/worker  += 40%
/home/user/projects/app         += 40%
/home/user/projects             += 40%
/home/user                      += 40%
/home                           += 40%
/                               += 40%
```

If a second process at `/home/user/projects/other` uses 20% CPU:

```
/home/user/projects/other       += 20%
/home/user/projects             += 20%  (now 60% total)
/home/user                      += 20%  (now 60% total)
/home                           += 20%  (now 60% total)
/                               += 20%  (now 60% total)
```

## The resulting view

After aggregating all processes:

```
   60.0%  /home/user/projects
   60.0%  /home/user
   60.0%  /home
   60.0%  /
   40.0%  /home/user/projects/app
   40.0%  /home/user/projects/app/worker
   20.0%  /home/user/projects/other
```

The root `/` accumulates the sum of all visible processes — it's always the largest number on a busy system.

## Why sums exceed 100%

This is expected and correct behavior. Each directory shows the **total CPU load of everything inside it**, not a normalized slice. The analogy is a budget tree where child costs roll up to parent totals.

If you see `/home/user/projects` at 350% on an 8-core machine, it means the processes in that subtree are collectively using 3.5 CPU cores.

## Data structure

The aggregation uses a plain `map[string]float64`:

```go
dirCPU := make(map[string]float64)
```

Keys are cleaned absolute paths. The map is rebuilt from scratch each cycle — no state persists between updates, which keeps memory bounded to the number of distinct directories across all running processes.

## Complexity

For each process, path traversal is `O(d)` where `d` is the directory depth. Total aggregation time is `O(n × d_avg)` where `n` is the number of processes above the threshold.

On a typical Linux desktop with ~300 processes and average depth ~5, this is roughly 1,500 string operations — fast enough to be imperceptible relative to the sampling interval.
