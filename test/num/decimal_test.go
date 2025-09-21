package num_test

import (
	"math"
	"testing"

	"cldr/num"
)

type decimalTestCase struct {
	locale   string
	whole    int64
	frac     uint64
	scale    int8
	expected string
}

func TestDecimalFormatter(t *testing.T) {
	t.Run("NewDecimalFormatter()", func(t *testing.T) {
		t.Run("unsupported locales", func(t *testing.T) {
			for i, tc := range []decimalTestCase{
				{"xx", 1000000, 100, 1, ""},
				{"en-XX", 1000000, 100, 1, ""},
			} {
				_, err := num.NewDecimalFormatter(tc.locale)
				if err == nil {
					t.Errorf("test case #%d - expected error but did not receive one", i+1)
				}
			}
		})
	})

	t.Run("Format()", func(t *testing.T) {
		t.Run("unsupported scale error", func(t *testing.T) {
			for i, tc := range []decimalTestCase{
				{"en", 1000000, 100, 21, ""},
				{"en", 1000000, 100, -2, ""},
			} {
				nf := num.MustNewDecimalFormatter(tc.locale)

				err := nf.SetScale(tc.scale)
				if err == nil {
					t.Errorf("test case #%d - expected error but did not receive one", i+1)
				}
			}
		})
		t.Run("fractional scale error", func(t *testing.T) {
			for i, tc := range []decimalTestCase{
				{"en", 1000000, 100, 1, ""},
				{"en", 1000000, math.MaxUint64, 19, ""},
			} {
				nf := num.MustNewDecimalFormatter(tc.locale)
				_ = nf.SetScale(tc.scale)

				_, err := nf.Format(tc.whole, tc.frac)
				if err == nil {
					t.Errorf("test case #%d - expected error but did not receive one", i+1)
				}
			}
		})

		t.Run("expected outputs", func(t *testing.T) {
			for i, tc := range []decimalTestCase{
				{"en", 1, 0, 0, "1"},
				{"en", 1, 0, 2, "1.00"},
				{"en", 1, 0o1, 0, "1"},
				{"en", 1, 0o1, -1, "1.1"},
				{"en", 1, 0o1, 1, "1.1"},
				{"en", 1, 0o1, 2, "1.01"},
				{"en", 1, 10, 2, "1.10"},
				{"en", 1, 10, -1, "1.1"},
				{"en", 1, 100, -1, "1.1"},
				{"en", 1, 101, -1, "1.101"},
				{"en", 1, 1010, -1, "1.101"},
				{"en", 1, 1010, 4, "1.1010"},
				{"en", 1, 1010, 5, "1.01010"},
				{"en", 1, 1010, 9, "1.000001010"},
				{"en", 1, math.MaxUint64, -1, "1.18446744073709551615"},
				{"en", 1, math.MaxUint64, 20, "1.18446744073709551615"},

				{"en-US", 100, 100, -1, "100.1"},
				{"en-US", 1000, 100, -1, "1,000.1"},
				{"en-US", 100000, 100, -1, "100,000.1"},
				{"en_US", 1000000, 100, -1, "1,000,000.1"},

				{"fr", 1000, 100, -1, "1\u202f000,1"},
				{"fr", 10000, 100, -1, "10\u202f000,1"},

				{"bn", 100000, 100, -1, "১,০০,০০০.১"},
				{"ar-YE", 100000, 100, -1, "١٠٠٬٠٠٠٫١"},
			} {
				nf := num.MustNewDecimalFormatter(tc.locale)
				_ = nf.SetLocale(tc.locale)
				nf.MustSetLocale(tc.locale)
				_ = nf.SetScale(tc.scale)
				nf.MustSetScale(tc.scale)

				actual, _ := nf.Format(tc.whole, tc.frac)
				if actual != tc.expected {
					t.Errorf("test case #%d - got: %v, expected: %v", i+1, actual, tc.expected)
				}

				actual = nf.MustFormat(tc.whole, tc.frac)
				if actual != tc.expected {
					t.Errorf("test case #%d - got: %v, expected: %v", i+1, actual, tc.expected)
				}
			}
		})
	})
}
