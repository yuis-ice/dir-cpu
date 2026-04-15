# Comparison: top / htop / glances

`dir-cpu` is not a replacement for these tools — it solves a different problem. Here's where each fits.

## Feature comparison

| Feature | `top` | `htop` | `glances` | `dir-cpu` |
|---------|-------|--------|-----------|-----------|
| Per-process CPU% | ✅ | ✅ | ✅ | ❌ |
| Per-directory CPU% | ❌ | ❌ | ❌ | ✅ |
| CPU by project tree | ❌ | ❌ | ❌ | ✅ |
| Memory usage | ✅ | ✅ | ✅ | ❌ |
| Network/disk I/O | ❌ | partial | ✅ | ❌ |
| Kill processes | ✅ | ✅ | ❌ | ❌ |
| Filter/search | partial | ✅ | ✅ | via grep |
| TUI / interactive | ✅ | ✅ | ✅ | ❌ |
| Works as a non-root | ✅ | ✅ | ✅ | partial |
| Linux only | ❌ | ❌ | ❌ | ✅ |

## When to use each

### Use `top` or `htop` when:

- You need to find and kill a specific runaway process
- You want to see memory, priority, or thread counts alongside CPU
- You want an interactive, keyboard-driven interface
- You're doing general-purpose system inspection

### Use `glances` when:

- You need a comprehensive system dashboard (CPU, memory, network, disk, temperatures)
- You're monitoring a server remotely and want everything in one screen
- You want alerting thresholds and logging

### Use `dir-cpu` when:

- You're a developer who wants to know **which project is eating CPU**, not which PID
- You're running multiple services in a monorepo and want per-service cost
- You see high CPU but can't tell which of your dozen terminal windows is responsible
- You want a quick forensics check for processes running from suspicious paths
- You want to answer "is my app or the system causing this heat?"

## The workflow difference

**With top/htop:**
```
1. Open htop
2. See PID 48291 at 85% CPU, command: "python3"
3. Wonder which python3 — your app? a cron job? conda?
4. Check: ls -la /proc/48291/cwd
5. Realize it's your ML training script
6. Close htop, go back to work
```

**With dir-cpu:**
```
1. Run dir-cpu
2. See /home/user/projects/ml-training at 85%
3. Done
```

## They complement each other

A common workflow:

1. **dir-cpu** to identify which project directory is the problem
2. **htop** filtered to that directory's processes (htop's `F4` filter by command) to inspect individual processes, memory, and thread counts

`dir-cpu` narrows the search space; `htop` lets you act on individual processes once you know where to look.
