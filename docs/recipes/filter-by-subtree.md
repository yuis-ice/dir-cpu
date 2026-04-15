# Filter by Subtree

`dir-cpu` shows the entire filesystem by default, which can be noisy on a busy system. Here are practical ways to narrow focus to a specific subtree.

## Using grep to filter output

Since `dir-cpu` writes to stdout and clears the screen with ANSI codes, direct piping doesn't work cleanly. Instead, use `watch` with `grep`:

```bash
watch -n 1 "dir-cpu -i 1s -n 200 -t 0 2>/dev/null | grep '^' | grep '/home/user/projects'"
```

Or strip the ANSI clear codes and pipe normally:

```bash
dir-cpu -i 1s -n 200 -t 0 | grep --line-buffered '/home/user/projects'
```

::: tip
The `-t 0` flag disables the threshold filter, ensuring all directories are present in the output before your `grep` applies its own filter.
:::

## Watching a single project

To watch only `/home/user/projects/myapp` and its subdirectories:

```bash
dir-cpu -t 0 -n 200 | grep --line-buffered 'myapp'
```

Sample output when only that subtree is active:

```
   82.3%  /home/user/projects/myapp
   82.3%  /home/user/projects/myapp/backend
   18.1%  /home/user/projects/myapp/worker
```

## Monitoring a monorepo

In a monorepo, several services run simultaneously under one root. Use `-t` to hide noise:

```bash
dir-cpu -base=cwd -t 1.0
```

This shows only directories where something meaningful is running, cutting the list from hundreds of entries to the handful of active services.

## One-shot snapshot

For a quick CPU attribution snapshot (not continuous), use `timeout`:

```bash
timeout 3s dir-cpu -i 2s -n 100 -t 0.1 2>/dev/null | tail -n +4
```

This captures one refresh cycle and exits.

## Script integration

To capture a snapshot to a file for later analysis:

```bash
#!/bin/bash
# Capture CPU-by-directory snapshot

OUTFILE="cpu-snapshot-$(date +%Y%m%d-%H%M%S).txt"

# Strip ANSI clear codes, capture after first refresh
timeout 3s dir-cpu -i 2s -n 200 -t 0 2>/dev/null \
  | sed 's/\x1b\[[0-9;]*[mJH]//g' \
  | grep -E '^\s+[0-9]' \
  > "$OUTFILE"

echo "Saved to $OUTFILE"
```
