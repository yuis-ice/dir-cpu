# Architecture Overview

## The directory-first paradigm

Traditional monitoring tools are **process-centric**: they show you a list of PIDs, command names, and their individual CPU percentages. This model maps directly to how the OS thinks about processes.

But developers don't think in PIDs. They think in projects.

When a developer asks "why is my laptop hot?", the answer they need is "your backend service is at 80% CPU" — not "PID 48291 (`node`) is at 80%". The second answer requires cross-referencing the PID against what you remember running, which directory it's in, and which project it belongs to.

`dir-cpu` inverts the model: instead of listing processes and asking you to mentally group them, it groups them for you by where they live on the filesystem.

## Data sources: /proc

Everything `dir-cpu` reads comes from Linux's `/proc` virtual filesystem:

| File | Used for |
|------|----------|
| `/proc/[pid]/stat` | CPU time (via `gopsutil`) |
| `/proc/[pid]/cwd` | Working directory (symlink) |
| `/proc/[pid]/exe` | Executable path (symlink) |

No sampling daemons, no eBPF, no kernel modules required.

## High-level flow

```
Snapshot 1                    Snapshot 2
─────────────────────────    ─────────────────────────
For each PID:                For each PID:
  read CPUTimes                read CPUTimes
  read cwd/exe path            read cwd/exe path
         │                            │
         └──────────── sleep ─────────┘
                            │
                    Compute delta per PID
                            │
                    Roll up through path hierarchy
                            │
                    Sort by total CPU%
                            │
                    Render to terminal
```

Two key design decisions fall out of this:

1. **One sleep per cycle** — the entire system is sampled before sleeping, not one process at a time. This makes cycle time `O(interval)` regardless of process count.

2. **Path hierarchy as the aggregation key** — every ancestor directory of a process's path gets that process's CPU added to it. A process at `/home/user/projects/app/worker` contributes to `worker`, `app`, `projects`, `home/user`, and ultimately `/`.

## Components

```
main.go
  ├── snapshot()        reads all /proc data into map[PID]procData
  ├── delta loop        computes CPU% from two snapshots
  ├── dir aggregation   rolls CPU up the path tree
  └── renderer          sorts and prints to terminal
```

See the detailed pages for each component:

- [CPU Calculation](./cpu-calculation) — the math behind the percentages
- [Tree Roll-up Algorithm](./tree-rollup) — how path hierarchy aggregation works
- [Performance & Overhead](./performance) — what `dir-cpu` itself costs
