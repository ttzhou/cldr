//go:build generate

package gen

import (
	"fmt"
	"log/slog"
)

// LocaleFiles generates files for all CLDR locales
// with coverage "modern".
func LocaleFiles(
	writeZipFileTo string,
	writeFilesToDir string,
	coverageLevel string,
) {
	cf, err := downloadAndOpen(writeZipFileTo)
	if err != nil {
		panic(err)
	}
	slog.Info(fmt.Sprintf("Parsing CLDR data from %s...", writeZipFileTo))
	cldrData := cf.getData()

	slog.Info(
		fmt.Sprintf(
			"Generating locale files with coverage %q in %s...",
			coverageLevel,
			writeFilesToDir,
		),
	)
	known, total, err := cldrData.writeLocaleFiles(writeFilesToDir, coverageLevel)
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
