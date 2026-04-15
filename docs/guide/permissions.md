# Permissions

## What a regular user can see

When running as a normal user, Linux restricts access to `/proc/[pid]/` for processes owned by other users. Specifically:

- `/proc/[pid]/cwd` — permission denied for other users' processes
- `/proc/[pid]/exe` — permission denied for other users' processes

`dir-cpu` silently skips any process it cannot read. This means:

- **You'll see your own processes accurately.**
- **System processes and other users' processes are invisible.**

This is usually fine for development use — you typically want to see where *your* work is spending CPU.

## Running with sudo for full system visibility

```bash
sudo dir-cpu
```

With root privileges, `dir-cpu` can read every process's `cwd` and `exe`, giving a complete picture of system-wide CPU consumption by directory.

::: warning
Running any monitoring tool as root has the usual caveats. `dir-cpu` does not write files or modify system state, but use `sudo` only when you need full visibility.
:::

## Why some directories still show 0%

Even with `sudo`, some processes may have no readable path:

- Kernel threads — they have no `exe` or meaningful `cwd`
- Processes in a different mount namespace or container
- Processes that exited between the two sampling snapshots

These are silently dropped from the output.

## Containerized environments

Inside a Docker container, `dir-cpu` sees only the processes in the container's PID namespace. `/proc` is scoped to the container, so paths reflect the container's filesystem, not the host's.

To monitor from the host, run `dir-cpu` on the host (with `sudo` if needed), not inside the container.
