package num

// A DecimalFormatter is a struct containing
// information necessary to format locale-aware
// decimal strings using CLDR data.
type DecimalFormatter struct {
	scale           int8
	numberFormatter numberFormatter
}

// NewDecimalFormatter returns a DecimalFormatter with
// no fixed scale (-1), set to locale `l`. A non-nil error is returned if
// the locale is not supported.
//
// See SetScale for a description of scales.
func NewDecimalFormatter(l string) (DecimalFormatter, error) {
	df := DecimalFormatter{}

	f, err := newNumberFormatter(l)
	if err != nil {
		return df, err
	}

	f.useStandardDecimalFormat()
	df.numberFormatter = f
	_ = df.SetScale(-1)

	return df, nil
}

// MustNewDecimalFormatter calls NewDecimalFormatter, and panics if its error result is not nil.
// Otherwise, it returns the non-error result.
func MustNewDecimalFormatter(l string) DecimalFormatter {
	df, err := NewDecimalFormatter(l)
	if err != nil {
		panic(err)
	}

	return df
}

// SetScale changes the scale considered when formatting.
// "scale" refers to the number of digits to the right of the decimal separator,
// though we permit the special case of -1 scale.
//
// A scale of -1 means we will format the decimal "regularly",
// e.g. whole 10, frac 1000 => 10.1000 => 10.1
//
// whereas a positive scale of 5 would be
// whole 10, frac 1000 => 10.01000
//
// An error is returned if the scale is unsupported, i.e.
// if it is <-1 or > the max supported scale = 20.
func (df *DecimalFormatter) SetScale(s int8) error {
	if s < -1 || int(s) > len(fracFormats)-1 {
		return unsupportedScaleError(s)
	}

	df.scale = s

	return nil
}

// MustSetScale calls SetScale, and panics if it returns an error.
func (df *DecimalFormatter) MustSetScale(s int8) {
	if err := df.SetScale(s); err != nil {
		panic(err.Error())
	}
}

// SetLocale changes the locale considered when formatting.
// An error is returned if the locale is not supported.
func (df *DecimalFormatter) SetLocale(l string) error {
	return df.numberFormatter.setLocale(l)
}

// MustSetLocale calls SetLocale, and panics if it returns an error.
func (df *DecimalFormatter) MustSetLocale(l string) {
	if err := df.SetLocale(l); err != nil {
		panic(err)
	}
}

// Format formats a given number's whole and fractional parts into a locale-aware string.
// A non-nil error is returned if the formatting cannot be done.
func (df DecimalFormatter) Format(w int64, f uint64) (string, error) {
	return df.numberFormatter.format(w, f, df.scale, "")
}

// MustFormat calls Format, and panics if there is an error.
func (df DecimalFormatter) MustFormat(w int64, f uint64) string {
	s, err := df.Format(w, f)
	if err != nil {
		panic(err)
	}

	return s
}
