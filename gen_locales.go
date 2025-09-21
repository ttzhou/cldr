//go:build generate

//go:generate go run -tags=generate gen_locales.go --coverage=modern
package main

import (
	"flag"

	"cldr/internal/gen"
)

func main() {
	coverageFlag := flag.String(
		"coverage",
		"modern",
		"Coverage level required to register locale (default 'modern')",
	)
	flag.Parse()

	gen.LocaleFiles(
		"./internal/resources/data/cldr-47.0.0.zip",
		"./internal/locale",
		*coverageFlag,
	)
}
