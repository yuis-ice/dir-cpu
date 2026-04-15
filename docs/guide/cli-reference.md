# CLI Reference

## Synopsis

```
dir-cpu [flags]
```

## Flags

### `-base string` <Badge type="info" text="default: cwd" />

Controls what filesystem path is used to represent each process.

| Value | Path source | Best for |
|-------|-------------|----------|
| `cwd` | `/proc/[pid]/cwd` — the process's working directory | Scripts, interpreted languages, anything run from a project folder |
| `exe` | `/proc/[pid]/exe` — the process's executable binary | Compiled programs, system daemons |

See [cwd vs exe mode](./cwd-vs-exe) for a deeper comparison.

---

### `-i duration` <Badge type="info" text="default: 1s" />

Sampling interval. Controls how often CPU usage is recalculated and the display refreshed.

Valid Go duration strings: `500ms`, `1s`, `2s`, `5s`, etc.

Shorter intervals give more responsive output but slightly higher overhead from `/proc` reads. Values below `200ms` are rarely useful since the delta will be very small.

```bash
dir-cpu -i 500ms   # fast refresh
dir-cpu -i 5s      # coarse but very low overhead
```

---

### `-t float` <Badge type="info" text="default: 0.5" />

Display threshold in percent. Directories whose total CPU usage is below this value are hidden.

Lowering the threshold reveals more entries, including light background tasks:

```bash
dir-cpu -t 0.0   # show everything including idle directories
dir-cpu -t 5.0   # show only directories consuming 5%+
```

---

### `-n int` <Badge type="info" text="default: 40" />

Maximum number of rows to print per cycle. When more directories exceed the threshold, a summary line is shown:

```
  ... 12 more
```

Increase this on tall terminals or when you need to see deep hierarchies:

```bash
dir-cpu -n 100
```

## Exit codes

| Code | Meaning |
|------|---------|
| `0`  | Clean exit (Ctrl+C) |
| `1`  | Fatal error reading `/proc` on startup |

## Notes

- Output is cleared each cycle with ANSI escape codes (`\033[2J\033[H`). Pipe-safe alternatives are not currently supported.
- CPU percentages may exceed 100% on multi-core systems. This is intentional and consistent with `top -H` behavior — each core contributes up to 100%.
