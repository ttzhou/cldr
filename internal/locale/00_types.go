// Package locale contains structures for managing localization specific information.
// This file is itself not generated, but contains "seed/source" types
// needed for running the "go generate" command for this directory, which will
// create the actual locale data files in the same directory.
package locale

import "strings"

type Locale struct {
	Code      string
	Language  string
	Territory string
	Variant   string
	Data      LocaleData
}

type LocaleData struct {
	NumberInfo          NumberInfo
	SupportedCurrencies map[string]CurrencyData
}

func (l Locale) Name() string {
	return strings.Trim(strings.Join([]string{l.Language, l.Territory, l.Variant}, "-"), "-")
}

type CurrencyData struct {
	MinorDigits         uint8
	DisplayCode         string
	DisplaySymbol       string
	DisplaySymbolNarrow string
}

type NumberFormat struct {
	PrimaryGroupSize   uint8
	SecondaryGroupSize uint8

	Prefix string
	Suffix string

	NegPrefix string
	NegSuffix string
}

type NumberFormats struct {
	StandardDecimal NumberFormat

	StandardCurrencySymbol   NumberFormat
	StandardCurrencyAlpha    NumberFormat
	StandardCurrencyNoSymbol NumberFormat

	AccountingCurrencySymbol   NumberFormat
	AccountingCurrencyAlpha    NumberFormat
	AccountingCurrencyNoSymbol NumberFormat
}

type NumberInfo struct {
	NumberSystem        string
	Digits              [10]string
	FractionalSeparator string
	GroupingSeparator   string

	Formats NumberFormats
}
