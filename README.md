# cldr

<!-- badges -->
[![Go Reference](https://pkg.go.dev/badge/github.com/ttzhou/cldr.svg)](https://pkg.go.dev/github.com/ttzhou/cldr)
![go](https://img.shields.io/github/go-mod/go-version/ttzhou/cldr)
[![codecov](https://codecov.io/gh/ttzhou/cldr/graph/badge.svg?token=SUU0ERUAST)](https://codecov.io/gh/ttzhou/cldr)
[![ci-checks](https://github.com/ttzhou/cldr/actions/workflows/ci.yml/badge.svg)](https://github.com/ttzhou/cldr/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ttzhou/cldr)](https://goreportcard.com/report/github.com/ttzhou/cldr)

# about

Module `cldr` contains various packages that leverage [Unicode CLDR](https://cldr.unicode.org/) data.

## packages

- `num`: utilities for (CLDR) locale-aware formatting of numerical amounts

## examples

### `num`

```go
package main

// go version 1.25.1

import (
	"fmt"

	"github.com/govalues/money"
	"github.com/ttzhou/cldr/num"
)

func main() {
	amt := money.MustNewAmount("USD", 91411206, 3)
    fmt.Println(amt) // USD 91411.206

	mf := num.MustNewMoneyFormatter("en-US")
	f, w, _ := amt.Int64(2)
    fmt.Println(f, w) // 91411 21
	fmt.Println(mf.MustFormat(f, uint64(w), amt.Curr().String())) // USD 91,411.21

	mf.DisplayCurrencyAsSymbolNarrow()
	fmt.Println(mf.MustFormat(f, uint64(w), amt.Curr().String())) // $91,411.21

    mf.MustSetLocale("fr-CA")
	fmt.Println(mf.MustFormat(f, uint64(w), amt.Curr().String())) // 91 411,21 $
	// fmt.Println(mf.MustFormat(f, uint64(w), "JPY"))
    // panic("in MoneyFormatter.MustFormat: fractional part 21 exceeds scale 0 (JPY)")
	fmt.Println(mf.MustFormat(f, 0, "JPY")) // 91 411 ¥
}
```

## why even build this

I wanted a toy project to learn golang, and have always found currency
localization interesting, so this seemed like a decent project.

There are a fair number of decimal/money packages already out there in the
golang community, but only one, [bojanz's currency package](github.com/bojanz/currency), 
seemed to also provide a robust localized formatting component. However, it tied it tightly
to its own implementation of money (using [cockroachDB's apd package](github.com/cockroachdb/apd)), 
while I preferred [govalues' money package](github.com/govalues/money) implementation. 
Unfortunately, `govalues` does not have localized formatting functionality, 
so I decided to build my own for my own purposes, and share it for anyone to use if they wish.

If it happens to be useful for even one person, I consider it a win.

# contributing

Development is currently only supported on POSIX-compliant environments.

Run `.dev/setup.sh` to setup local environment. The `Makefile` contains useful commands for dev related tasks.

Open PRs from your fork against this repository's `main` branch.

# changelog

No changelog will be maintained until this package reaches stable status 1.0.0. Expect the possibility of breaking changes at any point until then.

# TODOs

- [ ] benchmarking
- [ ] fuzzing

# license

MIT
