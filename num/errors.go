package num

import "fmt"

const (
	maxSupportedScale = uint8(len(fracFormats) - 1)
)

func unsupportedLocaleError(l string) error {
	return fmt.Errorf("unsupported locale: %q", l)
}

func unsupportedLocaleCurrencyError(c, l string) error {
	return fmt.Errorf("unsupported currency %q for locale %q", c, l)
}

func unsupportedScaleError(s int8) error {
	if s < -1 {
		return fmt.Errorf("scale %d must be at least -1", s)
	}

	return fmt.Errorf("scale %d exceeds max supported scale %d", s, maxSupportedScale)
}

func fractionalScaleError(f uint64, s uint8) error {
	return fmt.Errorf("fractional part %d exceeds scale %d", f, s)
}
