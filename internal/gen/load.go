//go:build generate

// Package gen contains methods for generating structs to work with Unicode CLDR data.
// These methods are not intended for external use and hence rarely
// perform error checking and are not thoroughly unit tested.
// Other packages do/should not reference this package directly.
package gen

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const (
	cldrVersion = "47.0.0"
	urlFormat   = "https://github.com/unicode-org/cldr-json/releases/download/%[1]s/cldr-%[1]s-json-full.zip"
)

type cldrZipFiles map[string]*zip.File

// DownloadAndOpen takes a path to a Unicode JSON CLDR zip file
// and returns a map of relevant files within the zip by filename.
func downloadAndOpen(to string) (cldrZipFiles, error) {
	to = filepath.Clean(to)
	if _, err := os.Open(to); errors.Is(err, os.ErrNotExist) {
		f, _ := os.Create(to)
		resp, _ := http.Get(fmt.Sprintf(urlFormat, cldrVersion))
		_, _ = io.Copy(f, resp.Body)
		_, _ = f.Close(), resp.Body.Close()
	}

	cldrFiles := make(map[string]*zip.File)

	za, err := zip.OpenReader(to)
	if err != nil {
		return cldrFiles, err
	}

	for _, f := range za.File {
		if isRelevantFile(f.Name) {
			cldrFiles[f.Name] = f
		}
	}

	return cldrFiles, err
}

func isRelevantFile(fn string) bool {
	for _, prefix := range []string{
		// currencies
		"cldr-bcp47/bcp47/currency.json",
		"cldr-core/supplemental/currencyData.json",

		// locales
		"cldr-core/availableLocales.json",
		"cldr-core/coverageLevels.json",
		"cldr-core/defaultContent.json",

		// locale naming
		"cldr-localenames-full/main/en/languages.json",
		"cldr-localenames-full/main/en/territories.json",
		"cldr-localenames-full/main/en/variants.json",

		// locale data
		"cldr-numbers-full/main/", //

		// number systems
		"cldr-core/supplemental/numberingSystems.json",
	} {
		if strings.HasPrefix(fn, prefix) {
			return true
		}
	}

	return false
}

type cldrData map[string]any

func (czf cldrZipFiles) getData() cldrData {
	locales := czf.getLocaleNamingInfo()
	localesData := czf.getLocalesData()
	localesMappings := czf.getLocaleMappings(locales, localesData)

	data := map[string]any{
		"locales":         locales,
		"locales-data":    localesData,
		"locale-mappings": localesMappings,
		"number-systems":  czf.getNumberingSystemsData(),
		"currencies":      czf.getCurrenciesData(),
	}

	return data
}

// All locale naming information known by CLDR; CLDR does not necessarily provide
// data for all of these, but we need them all to be able to assess
// valid locale name inputs.
type cldrLocaleNamingInfo map[string]map[string]string

func (czf cldrZipFiles) getLocaleNamingInfo() cldrLocaleNamingInfo {
	al, _ := czf["cldr-core/availableLocales.json"].Open()

	var availableLocales map[string]map[string][]string

	_ = json.NewDecoder(al).Decode(&availableLocales)
	_ = al.Close()

	dc, _ := czf["cldr-core/defaultContent.json"].Open()

	var defaultContent map[string][]string

	_ = json.NewDecoder(dc).Decode(&defaultContent)
	_ = dc.Close()

	acl := availableLocales["availableLocales"]["full"]
	dcl := defaultContent["defaultContent"]

	locales := slices.Concat(acl, dcl)
	slices.Sort(locales)

	llf, _ := czf["cldr-localenames-full/main/en/languages.json"].Open()

	var localeLanguagesDataRaw map[string]map[string]any

	_ = json.NewDecoder(llf).Decode(&localeLanguagesDataRaw)
	_ = llf.Close()

	localeLanguagesData := localeLanguagesDataRaw["main"]["en"].(map[string]any)["localeDisplayNames"].(map[string]any)["languages"].(map[string]any)

	ltf, _ := czf["cldr-localenames-full/main/en/territories.json"].Open()

	var localeTerritoriesDataRaw map[string]map[string]any

	_ = json.NewDecoder(ltf).Decode(&localeTerritoriesDataRaw)
	_ = ltf.Close()

	localeTerritoriesData := localeTerritoriesDataRaw["main"]["en"].(map[string]any)["localeDisplayNames"].(map[string]any)["territories"].(map[string]any)

	lvf, _ := czf["cldr-localenames-full/main/en/variants.json"].Open()

	var localeVariantsDataRaw map[string]map[string]any

	_ = json.NewDecoder(lvf).Decode(&localeVariantsDataRaw)
	_ = lvf.Close()

	localeVariantsData := localeVariantsDataRaw["main"]["en"].(map[string]any)["localeDisplayNames"].(map[string]any)["variants"].(map[string]any)

	clni := make(map[string]map[string]string)
	for _, code := range locales {
		clni[code] = make(map[string]string)
		lcs := strings.Split(code, "-")

		clni[code]["language"] = ""
		if language, ok := localeLanguagesData[lcs[0]]; ok {
			clni[code]["language"] = language.(string)
		}

		clni[code]["territory"] = ""

		if len(lcs) > 1 {
			if territory, ok := localeTerritoriesData[lcs[1]]; ok {
				clni[code]["territory"] = territory.(string)
			}
		}

		clni[code]["variant"] = ""

		if len(lcs) > 2 {
			if variant, ok := localeVariantsData[lcs[2]]; ok {
				clni[code]["variant"] = variant.(string)
			}
		}
	}

	return clni
}

// All locales agnostic currency data, e.g. currency minor units.
type cldrCurrenciesData map[string]map[string]string

func (czf cldrZipFiles) getCurrenciesData() cldrCurrenciesData {
	// Using this to get list of all known currency codes
	csdf, _ := czf["cldr-bcp47/bcp47/currency.json"].Open()

	var currencyFileData map[string]any

	_ = json.NewDecoder(csdf).Decode(&currencyFileData)

	currencyDescriptions := currencyFileData["keyword"].(map[string]any)["u"].(map[string]any)["cu"].(map[string]any)

	currencyCodes := make(map[string]string)

	for key, val := range currencyDescriptions {
		if strings.HasPrefix(key, "_") {
			continue
		}

		currencyCodes[strings.ToUpper(key)] = val.(map[string]any)["_description"].(string)
	}

	_ = csdf.Close()

	// Relevant currency data; doesn't contain all currencies,
	// instead using DEFAULT to fallback, so that is why we need
	// to get the list of all codes
	csdf, _ = czf["cldr-core/supplemental/currencyData.json"].Open()

	var currencySupplementalFileData map[string]map[string]map[string]any

	_ = json.NewDecoder(csdf).Decode(&currencySupplementalFileData)

	currenciesFractionData := currencySupplementalFileData["supplemental"]["currencyData"]["fractions"].(map[string]any)

	currenciesData := make(map[string]map[string]string)
	// If CLDR does not explicitly record currency data for a currency, use the DEFAULT data.
	defaultCurrencyFractionData := currenciesFractionData["DEFAULT"].(map[string]any)

	for code := range currencyCodes {
		currenciesData[code] = make(map[string]string)

		currencyFractionData := defaultCurrencyFractionData

		_, ok := currenciesFractionData[code]
		if ok {
			currencyFractionData = currenciesFractionData[code].(map[string]any)
		}

		// If non-zero, 10^digits = number of minor units in 1 major unit for currency
		// If 0, then the currency has no minor units.
		currenciesData[code]["digits"] = defaultCurrencyFractionData["_digits"].(string)

		digits, ok := currencyFractionData["_digits"]
		if ok {
			currenciesData[code]["digits"] = digits.(string)
		}
		// Not used currently in formatting/localization, but useful to parse
		currenciesData[code]["rounding"] = defaultCurrencyFractionData["_rounding"].(string)

		rounding, ok := currencyFractionData["_rounding"]
		if ok {
			currenciesData[code]["rounding"] = rounding.(string)
		}

		currenciesData[code]["cashDigits"] = currenciesData[code]["digits"]

		cashDigits, ok := currencyFractionData["_cashDigits"]
		if ok {
			currenciesData[code]["cashDigits"] = cashDigits.(string)
		}

		currenciesData[code]["cashRounding"] = "1"

		cashRounding, ok := currencyFractionData["_cashRounding"]
		if ok {
			currenciesData[code]["cashRounding"] = cashRounding.(string)
		}
	}

	_ = csdf.Close()

	return currenciesData
}

// Information about numbering systems; mainly for their digit runes.
type cldrNumberingSystemsData map[string]string

func (czf cldrZipFiles) getNumberingSystemsData() cldrNumberingSystemsData {
	lf, _ := czf["cldr-core/supplemental/numberingSystems.json"].Open()

	var fileMap map[string]map[string]any

	_ = json.NewDecoder(lf).Decode(&fileMap)

	numberingSystemsData := fileMap["supplemental"]["numberingSystems"].(map[string]any)
	numberingSystems := make(map[string]string)

	for system := range numberingSystemsData {
		data := numberingSystemsData[system].(map[string]any)
		if data["_type"].(string) == "algorithmic" {
			numberingSystems[system] = "algorithmic"

			continue
		}

		digits, ok := data["_digits"]
		if !ok {
			continue
		}

		numberingSystems[system] = digits.(string)
	}

	_ = lf.Close()

	return numberingSystems
}

// Actual locales data that is needed for our purposes, pulled from CLDR's recorded data.
// Note that this uses the CLDR JSON format, which I believe actually has some errors.
// We handle these in the parsing below in what I feel is a more sensical manner.
type cldrLocalesData map[string]map[string]any

func removeCurrencyPlaceholdersAndTrimSurroundingSpaces(p string) string {
	sb := strings.Builder{}
	sc := strings.Split(p, ";")
	for i, s := range sc {
		ss := strings.TrimSpace(strings.ReplaceAll(s, "¤", ""))
		sb.WriteString(ss)
		if i != len(sc)-1 {
			sb.WriteString(";")
		}
	}
	return sb.String()
}

func (czf cldrZipFiles) getLocalesData() cldrLocalesData {
	lf, _ := czf["cldr-core/availableLocales.json"].Open()

	var localeListData map[string]map[string][]string

	_ = json.NewDecoder(lf).Decode(&localeListData)

	_ = lf.Close()

	localesData := make(map[string]map[string]any)
	for _, locale := range localeListData["availableLocales"]["full"] {
		localesData[locale] = make(map[string]any)

		// Numbers
		localeNumbersFilename := fmt.Sprintf("cldr-numbers-full/main/%s/numbers.json", locale)

		var localeNumbersData map[string]map[string]map[string]any

		r, _ := czf[localeNumbersFilename].Open()
		_ = json.NewDecoder(r).Decode(&localeNumbersData)

		localeNumberDataFormats := localeNumbersData["main"][locale]["numbers"].(map[string]any)

		defaultNumberingSystem := localeNumberDataFormats["defaultNumberingSystem"].(string)
		nativeNumberingSystem := localeNumberDataFormats["otherNumberingSystems"].(map[string]any)["native"].(string)

		localesDataDecimalFormat := localeNumberDataFormats[fmt.Sprintf("decimalFormats-numberSystem-%s", defaultNumberingSystem)].(map[string]any)
		localesDataCurrencyFormat := localeNumberDataFormats[fmt.Sprintf("currencyFormats-numberSystem-%s", defaultNumberingSystem)].(map[string]any)

		localesData[locale]["numberSystem-default"] = defaultNumberingSystem
		localesData[locale]["numberSystem-native"] = nativeNumberingSystem

		localesData[locale]["standard-decimalFormat"] = localesDataDecimalFormat["standard"].(string)

		localesData[locale]["accounting-moneyFormat-symbol"] = localesDataCurrencyFormat["accounting"].(string)
		localesData[locale]["accounting-moneyFormat-alpha"] = localesData[locale]["accounting-moneyFormat-symbol"].(string)
		aantn, ok := localesDataCurrencyFormat["accounting-alphaNextToNumber"]
		if ok && strings.ContainsAny(aantn.(string), ";") {
			localesData[locale]["accounting-moneyFormat-alpha"] = aantn.(string)
		}

		localesData[locale]["accounting-moneyFormat-noSymbol"] = removeCurrencyPlaceholdersAndTrimSurroundingSpaces(
			localesDataCurrencyFormat["accounting"].(string),
		)
		anc, ok := localesDataCurrencyFormat["accounting-noCurrency"]
		if ok && strings.ContainsAny(anc.(string), ";") {
			localesData[locale]["accounting-moneyFormat-noSymbol"] = anc.(string)
		}

		localesData[locale]["standard-moneyFormat-symbol"] = localesDataCurrencyFormat["standard"]
		localesData[locale]["standard-moneyFormat-alpha"] = strings.Split(localesData[locale]["accounting-moneyFormat-alpha"].(string), ";")[0]
		localesData[locale]["standard-moneyFormat-noSymbol"] = removeCurrencyPlaceholdersAndTrimSurroundingSpaces(
			localesData[locale]["standard-moneyFormat-symbol"].(string),
		)

		// Separators, etc.
		localesSymbols := localeNumberDataFormats[fmt.Sprintf("symbols-numberSystem-%s", defaultNumberingSystem)].(map[string]any)

		for key, val := range localesSymbols {
			key = fmt.Sprintf("symbol-%s", key)
			localesData[locale][key] = val.(string)
		}

		_ = r.Close()

		// Decimal + currency/accounting symbols/formats
		lcf := fmt.Sprintf("cldr-numbers-full/main/%s/currencies.json", locale)

		var localeCurrenciesData map[string]map[string]map[string]any

		r, _ = czf[lcf].Open()
		_ = json.NewDecoder(r).Decode(&localeCurrenciesData)

		currencyFormats := make(map[string]map[string]string)

		localeCurrenciesNames := localeCurrenciesData["main"][locale]["numbers"].(map[string]any)
		for cur, dataRaw := range localeCurrenciesNames["currencies"].(map[string]any) {
			currencyFormats[cur] = make(map[string]string)
			data := dataRaw.(map[string]any)

			displayName, ok := data["displayName"]
			if ok {
				currencyFormats[cur]["display-name"] = displayName.(string)
			}

			symbol, ok := data["symbol"]
			if ok {
				currencyFormats[cur]["symbol"] = symbol.(string)
			}

			symbolNarrow, ok := data["symbol-alt-narrow"]
			if ok {
				currencyFormats[cur]["symbol-narrow"] = symbolNarrow.(string)
			}
		}

		localesData[locale]["currency-formats"] = currencyFormats
		_ = r.Close()
	}

	return localesData
}

// Synthetic - CLDR doesn't explicitly list all locale data for every locale,
// relying someimtes on "default" fallbacks. We construct this mapping so that
// we can still pull data on valid locales that do not necessarily have explicitly
// listed data, e.g. en-US (which maps to en, which does have data).
type cldrLocaleMappings map[string]map[string]string

func (czf cldrZipFiles) getLocaleMappings(
	ln cldrLocaleNamingInfo,
	ld cldrLocalesData,
) cldrLocaleMappings {
	cf, _ := czf["cldr-core/coverageLevels.json"].Open()

	var coverageLevelsRaw map[string]any

	_ = json.NewDecoder(cf).Decode(&coverageLevelsRaw)
	_ = cf.Close()

	coverageLevels := coverageLevelsRaw["effectiveCoverageLevels"].(map[string]any)
	localesMappings := make(map[string]map[string]string)

	for locale := range ln {
		lcs := strings.Split(locale, "-")
		for i := len(lcs) - 1; i >= 0; i-- {
			search := strings.Join(lcs[:i+1], "-")

			_, exists := ld[search]
			if exists {
				coverage, ok := coverageLevels[search]
				if !ok {
					coverage = "unknown"
				}

				localesMappings[locale] = map[string]string{
					"coverage":     coverage.(string),
					"known-locale": search,
				}

				break
			}
		}
	}

	return localesMappings
}
