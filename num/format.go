// Package num contains utilities for localized formatting
// of numerical amounts based on Unicode CLDR data.
package num

import (
	"fmt"
	"strconv"
	"strings"

	"cldr/internal/locale"
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
	localeCode string
	locale     locale.Locale

	numberFormat locale.NumberFormat
}

func newNumberFormatter(l string) (numberFormatter, error) {
	f := numberFormatter{}

	err := f.setLocale(l)
	if err != nil {
		return f, err
	}

	f.numberFormat = f.locale.Data.NumberInfo.Formats.StandardDecimal

	return f, nil
}

func (f numberFormatter) format(w int64, fn uint64, s int8, cs string) (string, error) {
	sb := strings.Builder{}

	isNegative := w < 0
	if isNegative {
		sb.WriteString(f.numberFormat.NegPrefix)

		w *= -1
	} else {
		sb.WriteString(f.numberFormat.Prefix)
	}

	sb.WriteString(f.formatWhole(uint64(w)))

	if s != 0 {
		fs, err := f.formatFrac(fn, uint8(max(s, 0)))
		if err != nil {
			return "", err
		}

		if len(fs) > 0 {
			sb.WriteString(f.locale.Data.NumberInfo.FractionalSeparator)
			sb.WriteString(fs)
		}
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

func (f numberFormatter) formatFrac(n uint64, mw uint8) (string, error) {
	var s string

	if mw > 0 {
		if mw > uint8(MaxSupportedScale) {
			return "", unsupportedScaleError(int8(mw))
		}
		if countDigits(n) > mw {
			return "", fractionalScaleError(n, mw)
		}

		i := min(mw, uint8(len(fracFormats)-1))
		s = fmt.Sprintf(fracFormats[i], n)
	} else {
		s = strings.TrimRight(strconv.FormatUint(n, 10), "0")
	}

	if f.locale.Data.NumberInfo.NumberSystem != "latn" {
		ss := strings.Split(s, "")
		for i, d := range ss {
			di, _ := strconv.Atoi(d)
			d = f.locale.Data.NumberInfo.Digits[di]
			ss[i] = d
		}

		s = strings.Join(ss, "")
	}

	return s, nil
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

	f.localeCode = l
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
