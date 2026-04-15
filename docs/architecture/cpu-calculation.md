# CPU Calculation

## Source data

CPU time for each process comes from `/proc/[pid]/stat`, read via `gopsutil`'s `p.Times()`. The relevant fields are:

- `utime` — CPU time spent in user mode (in clock ticks)
- `stime` — CPU time spent in kernel mode (in clock ticks)

`gopsutil` returns these as `float64` seconds via `cpu.TimesStat.User` and `cpu.TimesStat.System`.

## The snapshot-delta approach

`dir-cpu` takes **two snapshots** separated by the configured interval and computes the difference:

```
t₁ = snapshot before sleep
     time.Sleep(interval)
t₂ = snapshot after sleep
```

For each PID present in both snapshots:

```
delta = (t₂.User + t₂.System) - (t₁.User + t₁.System)
```

This delta is the number of CPU-seconds the process consumed during the interval.

## Converting to percentage

```
CPU% = (delta / interval_seconds) × 100
```

A process that consumed exactly `interval_seconds` worth of CPU time in one interval scores **100%** — meaning it fully utilized one CPU core for the entire interval.

### Multi-core behavior

On a system with N cores, percentages can reach `N × 100%`. A process using all 8 cores of an 8-core machine would show 800%.

This is consistent with how `top` reports CPU by default (per-core, not normalized). It is intentional — normalizing to the number of cores would make it impossible to tell whether a directory is using one core heavily or four cores lightly.

## Worked example

Interval: `1s`  
System: 4-core machine

```
PID 1234 — /home/user/projects/app/server.py
  t₁: User=10.500s, System=0.200s  → total = 10.700s
  t₂: User=11.100s, System=0.230s  → total = 11.330s

  delta = 11.330 - 10.700 = 0.630s
  CPU%  = (0.630 / 1.0) × 100 = 63.0%
```

This 63% then gets added to every ancestor directory of the process's cwd/exe path.

## Edge cases

### Negative delta

If a PID appears in both snapshots but shows lower cumulative CPU time in `t₂` (can happen after PID reuse or kernel counter wrap), the delta is clamped to zero:

```go
if delta <= 0 {
    continue
}
```

### PID not in previous snapshot

New processes that started between `t₁` and `t₂` are skipped — there's no baseline to compare against. They'll appear in the next cycle.

### PID disappeared before t₂

Processes that exited between snapshots are simply absent from `t₂` and thus not included in the output.

## Why not use gopsutil's CPUPercent(interval)?

`gopsutil`'s `p.CPUPercent(interval)` takes a duration argument and blocks for that duration **per call**. Calling it in a loop over all processes results in:

```
total_block_time ≈ N_processes × interval
```

On a system with 300 processes and a 1s interval, a single update cycle would take 5 minutes. `dir-cpu` avoids this entirely by taking bulk snapshots and computing deltas offline.
