# Security & Forensics

`dir-cpu` in `exe` mode has a practical security application: it surfaces processes running from unusual or suspicious filesystem locations. Legitimate software nearly always runs from well-known paths. Malicious processes often don't.

## The key insight

A cryptominer, backdoor, or persistence payload typically runs from:

- `/tmp/` — world-writable, common dropper target
- `/dev/shm/` — memory-mapped filesystem, sometimes used to avoid disk forensics
- `/var/tmp/` — less frequently cleaned than `/tmp`
- `/home/user/.local/share/` — hidden in dotdirs
- Random paths like `/opt/.hidden/`, `/.x/`, etc.

`dir-cpu -base=exe` immediately shows these because they appear as high-CPU directories in locations that have no business running executables.

## Quick scan for suspicious processes

```bash
sudo dir-cpu -base=exe -t 0.1 -i 2s
```

Look for any directory outside of:
- `/usr/bin`, `/usr/sbin`, `/usr/lib`, `/usr/local/`
- `/bin`, `/sbin`
- Known application paths like `/opt/known-app/`
- Your own project directories

Anything appearing under `/tmp`, `/dev/shm`, `/var/tmp`, or a dotdir that is consuming CPU is worth investigating immediately.

## One-shot suspicious path scan

```bash
#!/bin/bash
# Print any high-CPU exe paths that look suspicious

SUSPICIOUS_PATTERNS="/tmp/|/dev/shm/|/var/tmp/|\\.hidden|/proc/[0-9]"

timeout 4s sudo dir-cpu -base=exe -i 2s -t 0.5 -n 200 2>/dev/null \
  | sed 's/\x1b\[[0-9;]*[mJH]//g' \
  | grep -E '^\s+[0-9]' \
  | grep -E "$SUSPICIOUS_PATTERNS"
```

If this script prints anything, investigate those PIDs:

```bash
# Find PIDs running from /tmp
ls -la /proc/*/exe 2>/dev/null | grep '/tmp/'

# Get full info on a suspicious PID
cat /proc/<PID>/cmdline | tr '\0' ' '
ls -la /proc/<PID>/fd/
cat /proc/<PID>/maps | head -20
```

## Example: catching a cryptominer

A cryptominer dropped to `/tmp/kworker` and executed would show up as:

```
  380.0%  /tmp
  380.0%  /tmp/kworker     ← 380% on a 4-core machine = fully pegged
```

While legitimate kernel workers (`kworker`) run as kernel threads with no `exe` path and are invisible to `dir-cpu`, a binary *named* `kworker` dropped in `/tmp` would stand out immediately.

## Comparing against a baseline

On a known-clean system, capture a baseline:

```bash
timeout 4s sudo dir-cpu -base=exe -i 2s -t 0 -n 500 2>/dev/null \
  | sed 's/\x1b\[[0-9;]*[mJH]//g' \
  | grep -E '^\s+[0-9]' \
  | awk '{print $2}' \
  | sort > baseline-exe-dirs.txt
```

On a suspect system, compare:

```bash
timeout 4s sudo dir-cpu -base=exe -i 2s -t 0 -n 500 2>/dev/null \
  | sed 's/\x1b\[[0-9;]*[mJH]//g' \
  | grep -E '^\s+[0-9]' \
  | awk '{print $2}' \
  | sort > current-exe-dirs.txt

# Directories present now but not in baseline
comm -13 baseline-exe-dirs.txt current-exe-dirs.txt
```

New directories that appeared after a suspected compromise are immediate leads.

## Limitations

- `dir-cpu` only sees processes currently running and consuming CPU. A dormant backdoor waiting for a trigger won't appear.
- Processes that `exec` into a memfd (memory-only binary, no filesystem path) may show an empty `exe` path and be skipped.
- For deeper forensics, combine with `lsof`, `strace`, `auditd`, or tools like `osquery`.

`dir-cpu` is a first-pass triage tool, not a full forensics platform. Its value is speed: a suspicious directory surfaces in seconds with no setup.
