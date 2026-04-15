# FAQ

## Why does everything roll up to / at a high percentage?

This is correct behavior. `/` is the root of the entire filesystem, so it accumulates the CPU of every process on the system. It's always the highest number. Ignore `/` and focus on deeper paths.

## Why does the same percentage appear at multiple directory levels?

If only one process runs in a subtree, every ancestor from that process up to `/` shows the same number — because there's only one contributor being rolled up. When multiple processes in different subdirectories are running, you'll start seeing the numbers diverge at the branching point.

## My project isn't showing up

Check which mode you're using:

- **`-base=cwd`**: Your process needs to be run *from* your project directory. If it was launched from elsewhere (e.g., a system service started by init), `cwd` might be `/` or `/root`.
- **`-base=exe`**: Your binary needs to live inside your project directory or a known path.

Also verify the threshold: `-t 0.0` shows everything including idle processes.

## The percentages are way over 100%

That's expected on a multi-core machine. 100% means one full core. An 8-core machine can show up to 800% total. This is the same behavior as `top` in the default mode.

## Why does it need one interval before showing output?

The first cycle is a silent warmup: it takes the first snapshot and sleeps. Without this baseline, there's nothing to compare against and CPU% would be meaningless. After the first interval, output appears every cycle.

## Can I run it on macOS or Windows?

No. `dir-cpu` reads `/proc/[pid]/stat`, `/proc/[pid]/cwd`, and `/proc/[pid]/exe`, which are Linux-specific virtual filesystem paths. macOS and Windows have no `/proc`.

A macOS port would require using `proc_pidinfo` or `libproc`. Contributions welcome.

## Can I use it inside Docker?

Yes, but it will only see processes inside the container (scoped to the container's PID namespace). For host-level monitoring from inside a container you'd need `--pid=host`, which has security implications.

To monitor from the host, run `dir-cpu` directly on the host system.

## I see "(no significant usage)" even though my system is busy

All your processes are below the `-t` threshold (default 0.5%). Try:

```bash
dir-cpu -t 0.0
```

## Can I log the output to a file?

The display uses ANSI terminal escape codes to clear the screen each cycle, which makes direct redirection messy. For logging, strip the escape codes:

```bash
dir-cpu -i 5s | sed 's/\x1b\[[0-9;]*[mJH]//g' >> cpu-log.txt
```

Or use the snapshot script from the [Filter by Subtree](../recipes/filter-by-subtree) recipe.

## Why Go?

- Single static binary, no runtime dependencies
- Cross-compile easily for different Linux targets
- `gopsutil` provides clean access to `/proc` data with proper error handling
- Fast enough that overhead is negligible
