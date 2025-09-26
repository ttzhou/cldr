// Package num contains utilities for localized formatting
// of numerical amounts based on Unicode CLDR data.
package num

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ttzhou/cldr/internal/locale"
)

const (
	cldrCurSymbolPlaceholder = "¤"
)

// Up to 20 digits scale, which is the number of digits in largest uint64 number.
var fracFormats = [...]string{
	"%0d",
	"%01d",
	"%02d",
	"%03d",
	"%04d",
	"%05d",
	"%06d",
	"%07d",
	"%08d",
	"%09d",
	"%10d",
	"%11d",
	"%12d",
	"%13d",
	"%14d",
	"%15d",
	"%16d",
	"%17d",
	"%18d",
	"%19d",
	"%20d",
}

type numberFormatter struct {
	locale locale.Locale

	numberFormat locale.NumberFormat
}

func newNumberFormatter(l string) (numberFormatter, error) {
	f := numberFormatter{}

	err := f.setLocale(l)
	if err != nil {
		return f, err
	}
	f.useStandardDecimalFormat()

	return f, nil
}

func (f numberFormatter) format(w int64, fn uint64, s int8, cs string) (string, error) {
	fs, err := f.formatFrac(fn, s)
	if err != nil {
		return "", err
	}

	sb := strings.Builder{}

	isNegative := w < 0
	if isNegative {
		sb.WriteString(f.numberFormat.NegPrefix)

		w *= -1
	} else {
		sb.WriteString(f.numberFormat.Prefix)
	}

	sb.WriteString(f.formatWhole(uint64(w)))

	if len(fs) > 0 {
		sb.WriteString(f.locale.Data.NumberInfo.FractionalSeparator)
		sb.WriteString(fs)
	}

	if isNegative {
		sb.WriteString(f.numberFormat.NegSuffix)
	} else {
		sb.WriteString(f.numberFormat.Suffix)
	}

	ns := strings.NewReplacer(cldrCurSymbolPlaceholder, cs).Replace(sb.String())

	return ns, nil
}

func (f numberFormatter) formatWhole(n uint64) string {
	ns := strconv.FormatUint(n, 10)
	nss := strings.Split(ns, "")
	l, bufSize := len(nss), len(nss)

	pgs, sgs := int(f.numberFormat.PrimaryGroupSize), int(f.numberFormat.SecondaryGroupSize)
	if l > pgs {
		bufSize += ((l - 1 - pgs) / sgs) + 1
	}

	buf := make([]string, bufSize)

	gs, cnt, bi, ni := pgs, 0, bufSize-1, l-1
	for ni >= 0 {
		d := nss[ni]
		if f.locale.Data.NumberInfo.NumberSystem != "latn" {
			di, _ := strconv.Atoi(d)
			d = f.locale.Data.NumberInfo.Digits[di]
		}

		buf[bi] = d

		cnt++
		if cnt == gs && bi > 0 {
			bi--
			buf[bi] = f.locale.Data.NumberInfo.GroupingSeparator
			gs = sgs
			cnt = 0
		}

		bi--
		ni--
	}

	return strings.Join(buf, "")
}

func (f numberFormatter) formatFrac(n uint64, s int8) (string, error) {
	var ns string

	if s > 0 {
		us := uint8(s)
		if us > maxSupportedScale {
			return "", unsupportedScaleError(s)
		}
		if countDigits(n) > us {
			return "", fractionalScaleError(n, us)
		}

		ns = fmt.Sprintf(fracFormats[s], n)
	} else if s == 0 {
		if n > 0 {
			return "", fractionalScaleError(n, 0)
		}
		return "", nil
	} else if s == -1 {
		ns = strings.TrimRight(strconv.FormatUint(n, 10), "0")
	} else {
		return "", unsupportedScaleError(s)
	}

	if f.locale.Data.NumberInfo.NumberSystem != "latn" {
		ss := strings.Split(ns, "")
		for i, d := range ss {
			di, _ := strconv.Atoi(d)
			d = f.locale.Data.NumberInfo.Digits[di]
			ss[i] = d
		}

		ns = strings.Join(ss, "")
	}

	return ns, nil
}

func countDigits(n uint64) uint8 {
	if n == 0 {
		return 1
	}
	count := 0
	for n != 0 {
		n /= 10
		count++
	}
	return uint8(count)
}

func (f *numberFormatter) setLocale(l string) error {
	lc, ok := locale.Get(l)
	if !ok {
		return unsupportedLocaleError(l)
	}

	f.locale = lc

	return nil
}

func (f *numberFormatter) useStandardDecimalFormat() {
	f.numberFormat = f.locale.Data.NumberInfo.Formats.StandardDecimal
}

func (f *numberFormatter) useStandardCurrencyAlphaFormat() {
	f.numberFormat = f.locale.Data.NumberInfo.Formats.StandardCurrencyAlpha
}

func (f *numberFormatter) useStandardCurrencySymbolFormat() {
	f.numberFormat = f.locale.Data.NumberInfo.Formats.StandardCurrencySymbol
}

func (f *numberFormatter) useStandardCurrencyNoSymbolFormat() {
	f.numberFormat = f.locale.Data.NumberInfo.Formats.StandardCurrencyNoSymbol
}

func (f *numberFormatter) useAccountingCurrencyAlphaFormat() {
	f.numberFormat = f.locale.Data.NumberInfo.Formats.AccountingCurrencyAlpha
}

func (f *numberFormatter) useAccountingCurrencySymbolFormat() {
	f.numberFormat = f.locale.Data.NumberInfo.Formats.AccountingCurrencySymbol
}

func (f *numberFormatter) useAccountingCurrencyNoSymbolFormat() {
	f.numberFormat = f.locale.Data.NumberInfo.Formats.AccountingCurrencyNoSymbol
}
