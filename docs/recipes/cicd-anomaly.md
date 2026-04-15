# CI/CD Anomaly Detection

`dir-cpu` can be useful in CI/CD pipelines or build servers where you want to catch runaway processes, detect which build step is consuming unexpected CPU, or verify that a service started in the expected directory.

## Detecting runaway build processes

During a build, set a threshold and check whether anything unexpected is burning CPU outside the expected build directory:

```bash
#!/bin/bash
# Fail if any directory outside /opt/build consumes > 50% CPU

SNAPSHOT=$(timeout 4s dir-cpu -i 2s -n 200 -t 0 -base=cwd 2>/dev/null \
  | sed 's/\x1b\[[0-9;]*[mJH]//g' \
  | grep -E '^\s+[0-9]')

SUSPICIOUS=$(echo "$SNAPSHOT" | awk '$1+0 > 50 && $2 !~ /^\/opt\/build/')

if [ -n "$SUSPICIOUS" ]; then
  echo "WARNING: unexpected high-CPU directories:"
  echo "$SUSPICIOUS"
  exit 1
fi

echo "OK: CPU usage within expected directories"
```

## Verifying a service started in the right directory

After starting a background service, check that its working directory is where you expect:

```bash
# Start service
./myservice &
SVC_PID=$!
sleep 2

# Confirm it shows up under the expected directory in dir-cpu output
SNAPSHOT=$(timeout 4s dir-cpu -i 2s -n 200 -t 0 -base=cwd 2>/dev/null \
  | sed 's/\x1b\[[0-9;]*[mJH]//g')

if echo "$SNAPSHOT" | grep -q "/opt/myapp"; then
  echo "Service confirmed running under /opt/myapp"
else
  echo "ERROR: service not found under expected directory"
  kill $SVC_PID
  exit 1
fi
```

## Profiling a build step

Wrap a heavy build step with before/after snapshots to measure its CPU footprint by directory:

```bash
#!/bin/bash

capture_snapshot() {
  local label=$1
  timeout 4s dir-cpu -i 2s -n 200 -t 0.1 -base=cwd 2>/dev/null \
    | sed 's/\x1b\[[0-9;]*[mJH]//g' \
    | grep -E '^\s+[0-9]' \
    > "cpu-${label}.txt"
}

capture_snapshot "before"
make -j$(nproc) build
capture_snapshot "after"

echo "=== CPU during build (after snapshot) ==="
cat cpu-after.txt
```

## Using with process supervision

On a server running `supervisord` or `systemd`, you can use `dir-cpu` in `exe` mode to monitor service binaries:

```bash
# Watch only /opt/services/
watch -n 2 "timeout 3s dir-cpu -i 2s -base=exe -t 0 -n 100 2>/dev/null \
  | sed 's/\x1b\[[0-9;]*[mJH]//g' \
  | grep '/opt/services'"
```

Any service with unexpectedly high CPU will stand out immediately.
