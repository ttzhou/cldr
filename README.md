# cldr

<!-- badges -->
[![Go Reference](https://pkg.go.dev/badge/github.com/ttzhou/cldr.svg)](https://pkg.go.dev/github.com/ttzhou/cldr)
![go](https://img.shields.io/github/go-mod/go-version/ttzhou/cldr)
[![codecov](https://codecov.io/gh/ttzhou/cldr/graph/badge.svg?token=SUU0ERUAST)](https://codecov.io/gh/ttzhou/cldr)

# about

Module `cldr` contains various packages that leverage [Unicode CLDR](https://cldr.unicode.org/) data.

## packages

- `num`: utilities for (CLDR) locale-aware formatting of numerical amounts

# contributing

Development is currently only supported on POSIX-compliant environments.

Run `.dev/setup.sh` to setup local environment. The `Makefile` contains useful commands for dev related tasks.

Open PRs from your fork against this repository's `main` branch.

# TODOs

- [ ] benchmarking
- [ ] fuzzing

# license

MIT
