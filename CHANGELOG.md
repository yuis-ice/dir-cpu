# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-04-15

### Added

- Real-time CPU usage aggregated by filesystem directory
- `cwd` mode: attribute CPU to each process's working directory
- `exe` mode: attribute CPU to the directory containing each process's binary
- Snapshot-delta sampling — one sleep per cycle regardless of process count
- Tree roll-up: every ancestor directory accumulates child process CPU costs
- `-i` flag: configurable update interval (default 1s)
- `-t` flag: display threshold to filter low-CPU directories (default 0.5%)
- `-n` flag: max rows to display (default 40)
- Graceful handling of permission-denied processes (silently skipped)

[0.1.0]: https://github.com/yuis-ice/dir-cpu/releases/tag/v0.1.0
