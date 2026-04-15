package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/process"
)

type procData struct {
	times *cpu.TimesStat
	path  string
}

func snapshot(base string) map[int32]procData {
	procs, err := process.Processes()
	if err != nil {
		return nil
	}

	data := make(map[int32]procData, len(procs))
	for _, p := range procs {
		ts, err := p.Times()
		if err != nil {
			continue
		}

		var path string
		if base == "exe" {
			path, _ = p.Exe()
		} else {
			path, _ = p.Cwd()
		}

		data[p.Pid] = procData{times: ts, path: path}
	}
	return data
}

func main() {
	base := flag.String("base", "cwd", "aggregation basis: cwd | exe")
	interval := flag.Duration("i", 1*time.Second, "update interval")
	thresh := flag.Float64("t", 0.5, "display threshold (%)")
	maxRows := flag.Int("n", 40, "max rows to display")
	flag.Parse()

	fmt.Printf("dir-cpu (base=%s, interval=%v) — Ctrl+C to quit\n\n", *base, *interval)

	prev := snapshot(*base)
	if prev == nil {
		fmt.Fprintln(os.Stderr, "error: failed to read processes")
		os.Exit(1)
	}
	time.Sleep(*interval)

	for {
		curr := snapshot(*base)
		if curr == nil {
			time.Sleep(*interval)
			continue
		}

		secs := interval.Seconds()
		dirCPU := make(map[string]float64)

		for pid, c := range curr {
			p, ok := prev[pid]
			if !ok || c.path == "" || c.path == "/" {
				continue
			}

			delta := (c.times.User + c.times.System) - (p.times.User + p.times.System)
			if delta <= 0 {
				continue
			}

			cpuPct := delta / secs * 100
			if cpuPct < *thresh {
				continue
			}

			dir := filepath.Clean(c.path)
			for {
				dirCPU[dir] += cpuPct
				if dir == "/" {
					break
				}
				parent := filepath.Dir(dir)
				if parent == dir {
					break
				}
				dir = parent
			}
		}

		type entry struct {
			dir string
			pct float64
		}
		list := make([]entry, 0, len(dirCPU))
		for d, p := range dirCPU {
			list = append(list, entry{d, p})
		}
		sort.Slice(list, func(i, j int) bool { return list[i].pct > list[j].pct })

		fmt.Printf("\033[2J\033[H")
		fmt.Printf("CPU%% by %s dir  (updated: %s)\n", *base, time.Now().Format("15:04:05"))
		fmt.Println("────────────────────────────────────────")
		if len(list) == 0 {
			fmt.Println("(no significant usage)")
		}
		for i, e := range list {
			if i >= *maxRows {
				fmt.Printf("  ... %d more\n", len(list)-i)
				break
			}
			fmt.Printf("%7.1f%%  %s\n", e.pct, e.dir)
		}

		prev = curr
		time.Sleep(*interval)
	}
}
