---
layout: home

hero:
  name: dir-cpu
  text: CPU usage by directory
  tagline: See which project folder is burning CPU — not just which process.
  actions:
    - theme: brand
      text: Get Started
      link: /guide/getting-started
    - theme: alt
      text: How it works
      link: /architecture/overview
    - theme: alt
      text: GitHub
      link: https://github.com/yuis-ice/dir-cpu

features:
  - icon: 📂
    title: Directory-first view
    details: Rolls every process's CPU cost up through its ancestor directories. Spot expensive project trees at a glance without hunting through process lists.
  - icon: ⚡
    title: Accurate, low-overhead sampling
    details: Takes two /proc snapshots and computes the delta — one sleep per cycle regardless of how many processes are running.
  - icon: 🔍
    title: cwd or exe mode
    details: Group by working directory (ideal for Python/Node scripts) or by executable path (ideal for compiled binaries). Switch with a single flag.
  - icon: 🛡️
    title: Security & forensics
    details: Processes hiding in /tmp or unusual paths surface immediately in exe mode — useful for spotting rogue miners or backdoors.
---
