//go:build generate

package gen

import (
	"fmt"
	"log/slog"
	"path/filepath"
)

// LocaleFiles generates files for all CLDR locales
// with coverage "modern".
func LocaleFiles(
	dataDir string,
	localeFileDir string,
	coverageLevel string,
) {
	filename := fmt.Sprintf("cldr-%s.zip", cldrVersion)
	cf, err := downloadAndOpen(filepath.Join(dataDir, filename))
	if err != nil {
		panic(err)
	}
	slog.Info(fmt.Sprintf("Parsing CLDR data from %s...", dataDir))
	cldrData := cf.getData()

	slog.Info(
		fmt.Sprintf(
			"Generating locale files with coverage %q in %s...",
			coverageLevel,
			localeFileDir,
		),
	)
	known, total, err := cldrData.writeLocaleFiles(localeFileDir, coverageLevel)
	if err != nil {
		panic(err)
	}
	slog.Info(fmt.Sprintf(
		"Wrote %d CLDR locale files with coverage %q out of %d total locales...",
		known,
		coverageLevel,
		total,
	))
	slog.Info("Done!")
}
