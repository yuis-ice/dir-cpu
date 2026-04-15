# Contributing

Bug reports and pull requests are welcome.

## Getting started

```bash
git clone https://github.com/yuis-ice/dir-cpu
cd dir-cpu
go build .
```

## Guidelines

- Keep changes focused. One fix or feature per PR.
- Test on Linux before submitting (the tool reads `/proc` and won't work elsewhere).
- Match the existing code style — `gofmt` your changes.
- If adding a flag or changing behavior, update `README.md` accordingly.

## Reporting bugs

Open a GitHub issue with:
- What you ran (exact command)
- What you expected
- What actually happened
- OS and Go version (`go version`)

## Questions

Open an issue with the `question` label.
