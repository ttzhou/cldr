//go:build generate

package gen

// This file generates types that transform CLDR data
// into usable formats that can be written to generated code files.
// See the 00_types.go file in the internal `locale` package.
// As with other functions in this internal `gen` package, efficiency
// was not a primary concern.
import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"cldr/internal/locale"
)

// Parse a number format string into useful structured information.
// https://cldr.unicode.org/translation/number-currency-formats/number-and-currency-patterns
func generateNumberFormat(numfmtstr string) locale.NumberFormat {
	numfmt := locale.NumberFormat{}

	pgsregex := regexp.MustCompile(`,#+0.`)

	pgs := len(pgsregex.FindString(numfmtstr)) - 2
	if pgs <= 0 {
		pgs = 3
	}

	numfmt.PrimaryGroupSize = uint8(pgs)

	sgsregex := regexp.MustCompile(`,#+,`)

	sgs := len(sgsregex.FindString(numfmtstr)) - 2
	if sgs <= 0 {
		sgs = 3
	}

	numfmt.SecondaryGroupSize = uint8(sgs)

	numregex := regexp.MustCompile(`[0.,#]+`)
	numfmtstr = numregex.ReplaceAllString(numfmtstr, "0")

	nfsc := strings.Split(numfmtstr, ";")

	pfs := nfsc[0]
	posFixes := strings.Split(pfs, "0")

	var negFixes []string

	if len(nfsc) > 1 {
		nfs := nfsc[len(nfsc)-1]
		negFixes = strings.Split(nfs, "0")
	} else {
		negFixes = slices.Clone(posFixes)
	}
	// Explicitly set a style for the negative prefix case.
	// Negative sign always goes on the right of the prefix.
	if negFixes[0] == posFixes[0] {
		negFixes[0] = posFixes[0] + "-"
	}

	numfmt.Prefix = posFixes[0]
	numfmt.Suffix = posFixes[len(posFixes)-1]
	numfmt.NegPrefix = negFixes[0]
	numfmt.NegSuffix = negFixes[len(negFixes)-1]

	return numfmt
}

func (c cldrData) generateNumberInfo(l string) (locale.NumberInfo, error) {
	var nf locale.NumberInfo

	localesData := c["locales-data"].(cldrLocalesData)

	localeData, ok := localesData[l]
	if !ok {
		return nf, fmt.Errorf("locale %s does not exist", l)
	}

	defaultNumberSystem := localeData["numberSystem-default"].(string)

	nf.NumberSystem = defaultNumberSystem
	nf.Digits = [10]string(strings.Split(
		c["number-systems"].(cldrNumberingSystemsData)[localeData["numberSystem-default"].(string)],
		"",
	))
	nf.FractionalSeparator = localeData["symbol-decimal"].(string)
	nf.GroupingSeparator = localeData["symbol-group"].(string)
	nf.Formats = locale.NumberFormats{
		StandardDecimal: generateNumberFormat(
			localeData["standard-decimalFormat"].(string),
		),
		StandardCurrencySymbol: generateNumberFormat(
			localeData["standard-moneyFormat-symbol"].(string),
		),
		StandardCurrencyAlpha: generateNumberFormat(
			localeData["standard-moneyFormat-alpha"].(string),
		),
		StandardCurrencyNoSymbol: generateNumberFormat(
			localeData["standard-moneyFormat-noSymbol"].(string),
		),

		AccountingCurrencySymbol: generateNumberFormat(
			localeData["accounting-moneyFormat-symbol"].(string),
		),
		AccountingCurrencyAlpha: generateNumberFormat(
			localeData["accounting-moneyFormat-alpha"].(string),
		),
		AccountingCurrencyNoSymbol: generateNumberFormat(
			localeData["accounting-moneyFormat-noSymbol"].(string),
		),
	}

	return nf, nil
}

func (c cldrData) generateCurrencyData(l, cur string) (locale.CurrencyData, error) {
	var cd locale.CurrencyData

	localedata, ok := c["locales-data"].(cldrLocalesData)[l]
	if !ok {
		return cd, fmt.Errorf("locale %s does not exist", l)
	}

	currencydata, ok := c["currencies"].(cldrCurrenciesData)[cur]
	if !ok {
		return cd, fmt.Errorf("currency %s does not exist", cur)
	}

	ds, _ := strconv.Atoi(currencydata["digits"])
	cd.MinorDigits = uint8(ds)

	currencyFormats := localedata["currency-formats"].(map[string]map[string]string)

	currencyFormat, ok := currencyFormats[cur]
	if !ok {
		return cd, fmt.Errorf("currency %s does not exist", cur)
	}

	cd.DisplayCode = cur
	cd.DisplaySymbol = cur
	cd.DisplaySymbolNarrow = cur

	symbol, ok := currencyFormat["symbol"]
	if ok {
		cd.DisplaySymbol = symbol
		cd.DisplaySymbolNarrow = symbol
	}

	symbolNarrow, ok := currencyFormat["symbol-narrow"]
	if ok {
		cd.DisplaySymbolNarrow = symbolNarrow
	}

	return cd, nil
}

func (c cldrData) GenerateLocaleData(l string) (locale.LocaleData, error) {
	var ld locale.LocaleData

	localedata, ok := c["locales-data"].(cldrLocalesData)[l]
	if !ok {
		return ld, fmt.Errorf("locale %s does not exist", l)
	}

	numberInfo, err := c.generateNumberInfo(l)
	if !ok {
		return ld, err
	}

	ld.NumberInfo = numberInfo

	currenciesMap := make(map[string]locale.CurrencyData)
	for cur := range localedata["currency-formats"].(map[string]map[string]string) {
		currency, err := c.generateCurrencyData(l, cur)
		if err != nil {
			return ld, err
		}

		currenciesMap[cur] = currency
	}

	ld.SupportedCurrencies = currenciesMap

	return ld, nil
}

func (c cldrData) GenerateLocale(l string) (locale.Locale, error) {
	var lc locale.Locale

	locale, ok := c["locales"].(cldrLocaleNamingInfo)[l]
	if !ok {
		return lc, fmt.Errorf("locale %s does not exist", l)
	}

	lc.Code = l
	lc.Language = locale["language"]
	lc.Territory = locale["territory"]
	lc.Variant = locale["variant"]

	return lc, nil
}
