package num

import (
	"fmt"
	"slices"
	"unicode"

	"github.com/ttzhou/cldr/internal/locale"
)

// A MoneyFormatter can be used to format locale-aware
// monetary amounts using CLDR data.
type MoneyFormatter struct {
	useAccountingStyle bool

	currencyStyle currencyStyle
	currencyLabel currencyLabel

	numberFormatter numberFormatter
}

// NewMoneyFormatter returns a [MoneyFormatter] with
// no fixed scale (-1), set to locale `l`.
// A non-nil error is returned ift he locale is not supported.
func NewMoneyFormatter(l string) (MoneyFormatter, error) {
	mf := MoneyFormatter{}

	f, err := newNumberFormatter(l)
	if err != nil {
		return mf, err
	}

	mf.numberFormatter = f
	mf.UseStandardStyle()
	mf.DisplayCurrencyAsCode()

	return mf, nil
}

// MustNewMoneyFormatter calls [NewMoneyFormatter], and panics if its error result is not nil.
// Otherwise, it returns the non-error result.
func MustNewMoneyFormatter(l string) MoneyFormatter {
	mf, err := NewMoneyFormatter(l)
	if err != nil {
		panic(err)
	}

	return mf
}

// SetLocale changes the locale considered when formatting.
// An error is returned if the locale is not supported.
func (mf *MoneyFormatter) SetLocale(l string) error {
	return mf.numberFormatter.setLocale(l)
}

// MustSetLocale calls [MoneyFormatter.SetLocale], and panics if it returns an error.
func (mf *MoneyFormatter) MustSetLocale(l string) {
	if err := mf.SetLocale(l); err != nil {
		panic(err)
	}
}

// UseStandardStyle indicates that monetary amounts should be formatted
// in the standard, non-accounting style defined by CLDR for the current locale, if relevant.
func (mf *MoneyFormatter) UseStandardStyle() {
	mf.useAccountingStyle = false
}

// UseAccountingStyle indicates that monetary amounts should be formatted
// in the accounting style defined by CLDR for the current locale, if relevant.
func (mf *MoneyFormatter) UseAccountingStyle() {
	mf.useAccountingStyle = true
}

// DisplayCurrencyAsCode informs the formatter to format currency labels as its 3 letter ISO code.
func (mf *MoneyFormatter) DisplayCurrencyAsCode() {
	mf.currencyStyle = code
}

// DisplayCurrencyAsSymbol informs the formatter to format currency labels as its CLDR symbol.
func (mf *MoneyFormatter) DisplayCurrencyAsSymbol() {
	mf.currencyStyle = symbol
}

// DisplayCurrencyAsSymbolNarrow informs the formatter to format currency labels as its CLDR symbol, narrow variant.
func (mf *MoneyFormatter) DisplayCurrencyAsSymbolNarrow() {
	mf.currencyStyle = symbolnarrow
}

// DisplayNoCurrency informs the formatter to format with no currency symbol.
func (mf *MoneyFormatter) DisplayNoCurrency() {
	mf.currencyStyle = none
}

// Format formats a given number's whole and fractional parts into a locale-aware string
// for the given currency.
// A non-nil error is returned if:
//   - the currency is not supported for the formatter's currently set locale
//   - the fractional part exceeds the number of minor digits for the currency (e.g. 2 minor units would fail for JPY, which has no minor, but not for USD)
func (mf MoneyFormatter) Format(w int64, f uint64, c string) (string, error) {
	ci, setCurrencyErr := mf.setCurrency(c)
	if setCurrencyErr != nil {
		return "", setCurrencyErr
	}

	s, formatErr := mf.numberFormatter.format(w, f, int8(ci.MinorDigits), string(mf.currencyLabel))
	if formatErr != nil {
		return "", fmt.Errorf("%w (%s)", formatErr, c)
	}
	return s, nil
}

// MustFormat calls [MoneyFormatter.Format], and panics if it returns a non-nil error.
func (mf MoneyFormatter) MustFormat(w int64, f uint64, c string) string {
	s, err := mf.Format(w, f, c)
	if err != nil {
		panic(err.Error())
	}

	return s
}

type currencyStyle uint8

const (
	code currencyStyle = iota
	symbol
	symbolnarrow
	none
)

type currencyLabel string

func (cl currencyLabel) isEmpty() bool {
	return len(cl) == 0
}

func (cl currencyLabel) containsAlphaChars() bool {
	return slices.ContainsFunc([]rune(cl), func(r rune) bool { return unicode.IsLetter(r) })
}

func (mf *MoneyFormatter) setCurrency(c string) (locale.CurrencyData, error) {
	cd, ok := mf.numberFormatter.locale.Data.SupportedCurrencies[c]
	if !ok {
		return cd, unsupportedLocaleCurrencyError(c, mf.numberFormatter.locale.Code)
	}

	switch mf.currencyStyle {
	case code:
		mf.currencyLabel = currencyLabel(cd.DisplayCode)
	case symbol:
		mf.currencyLabel = currencyLabel(cd.DisplaySymbol)
	case symbolnarrow:
		mf.currencyLabel = currencyLabel(cd.DisplaySymbolNarrow)
	case none:
		mf.currencyLabel = currencyLabel("")
	}

	if !mf.useAccountingStyle {
		if mf.currencyLabel.isEmpty() {
			mf.numberFormatter.useStandardCurrencyNoSymbolFormat()
		} else if mf.currencyLabel.containsAlphaChars() {
			mf.numberFormatter.useStandardCurrencyAlphaFormat()
		} else {
			mf.numberFormatter.useStandardCurrencySymbolFormat()
		}
	} else {
		if mf.currencyLabel.isEmpty() {
			mf.numberFormatter.useAccountingCurrencyNoSymbolFormat()
		} else if mf.currencyLabel.containsAlphaChars() {
			mf.numberFormatter.useAccountingCurrencyAlphaFormat()
		} else {
			mf.numberFormatter.useAccountingCurrencySymbolFormat()
		}
	}

	return cd, nil
}
