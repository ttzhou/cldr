package num_test

import (
	"testing"

	"github.com/ttzhou/cldr/num"
)

type moneyTestCase struct {
	locale   string
	whole    int64
	frac     uint64
	cur      string
	expected string
}

func TestMoneyFormatter(t *testing.T) {
	t.Run("NewMoneyFormatter()", func(t *testing.T) {
		t.Run("unsupported locales", func(t *testing.T) {
			for i, tc := range []moneyTestCase{
				{"xx", 1000000, 100, "USD", "unsupported locale: \"xx\""},
				{"en-XX", 1000000, 100, "USD", "unsupported locale: \"en-XX\""},
			} {
				_, err := num.NewMoneyFormatter(tc.locale)
				if err == nil {
					t.Errorf("test case #%d - expected error but did not receive one", i+1)
					continue
				}
				got := err.Error()
				if got != tc.expected {
					t.Errorf("test case #%d - got: %s, expected: %s", i+1, got, tc.expected)
				}
			}
		})
	})

	t.Run("Format()", func(t *testing.T) {
		t.Run("fractional scale error", func(t *testing.T) {
			for i, tc := range []moneyTestCase{
				{"en", 1000000, 2, "JPY", "fractional part 2 exceeds scale 0 (JPY)"},
				{"en", 1000000, 240, "USD", "fractional part 240 exceeds scale 2 (USD)"},
			} {
				mf := num.MustNewMoneyFormatter(tc.locale)
				_, err := mf.Format(tc.whole, tc.frac, tc.cur)
				if err == nil {
					t.Errorf("test case #%d - expected error but did not receive one", i+1)
					continue
				}
				got := err.Error()
				if got != tc.expected {
					t.Errorf("test case #%d - got: %s, expected: %s", i+1, got, tc.expected)
				}
			}
		})

		t.Run("unsupported currencies for locale", func(t *testing.T) {
			for i, tc := range []moneyTestCase{
				{"en", 100000, 1, "XYX", "unsupported currency \"XYX\" for locale \"en\""},
			} {
				mf := num.MustNewMoneyFormatter(tc.locale)
				mf.UseStandardStyle()
				mf.DisplayCurrencyAsCode()

				_, err := mf.Format(tc.whole, tc.frac, tc.cur)
				if err == nil {
					t.Errorf("test case #%d - expected error but did not receive one", i+1)
					continue
				}
				got := err.Error()
				if got != tc.expected {
					t.Errorf("test case #%d - got: %s, expected: %s", i+1, got, tc.expected)
				}
			}
		})

		t.Run("standard style, currency code", func(t *testing.T) {
			for i, tc := range []moneyTestCase{
				{"en", 1, 0, "USD", "USD\u00a01.00"},
				{"en", 1, 2, "USD", "USD\u00a01.02"},
				{"en", 1, 0o1, "USD", "USD\u00a01.01"},

				{"fr", 1000, 10, "USD", "1\u202f000,10\u00a0USD"},
				{"fr", 10000, 10, "USD", "10\u202f000,10\u00a0USD"},
				{"fr", -10000, 10, "USD", "-10\u202f000,10\u00a0USD"},

				{"bn", 100000, 1, "USD", "১,০০,০০০.০১\u00a0USD"},
				{"bn", 100000, 10, "USD", "১,০০,০০০.১০\u00a0USD"},

				{"en", 100000, 1, "BHD", "BHD\u00a0100,000.001"},
				{"ar", 100000, 1, "BHD", "\u061c100,000.001\u00a0BHD"},
			} {
				mf := num.MustNewMoneyFormatter(tc.locale)
				_ = mf.SetLocale(tc.locale)
				mf.MustSetLocale(tc.locale)
				mf.UseStandardStyle()
				mf.DisplayCurrencyAsCode()

				actual := mf.MustFormat(tc.whole, tc.frac, tc.cur)
				if actual != tc.expected {
					t.Errorf("test case #%d - got: %v, expected: %v", i+1, actual, tc.expected)
				}
			}
		})

		t.Run("standard style, currency symbol", func(t *testing.T) {
			for i, tc := range []moneyTestCase{
				{"en", 1, 0o1, "USD", "$1.01"},
				{"fr", 1000, 10, "USD", "1\u202f000,10\u00a0$US"},
				{"fr-CA", 1000, 10, "USD", "1\u00a0000,10\u00a0$\u00a0US"},
				{"en-CA", 1000, 10, "USD", "US$\u00a01,000.10"},
				{"bn", 100000, 1, "USD", "১,০০,০০০.০১\u00a0US$"},
				{"en", 100000, 1, "BHD", "BHD\u00a0100,000.001"},
				{"ar", 100000, 1, "BHD", "\u061c100,000.001\u00a0د.ب.\u200f"},
			} {
				mf := num.MustNewMoneyFormatter(tc.locale)
				_ = mf.SetLocale(tc.locale)
				mf.MustSetLocale(tc.locale)
				mf.UseStandardStyle()
				mf.DisplayCurrencyAsSymbol()

				actual := mf.MustFormat(tc.whole, tc.frac, tc.cur)
				if actual != tc.expected {
					t.Errorf("test case #%d - got: %v, expected: %v", i+1, actual, tc.expected)
				}
			}
		})

		t.Run("standard style, currency symbol narrow", func(t *testing.T) {
			for i, tc := range []moneyTestCase{
				{"en", 1, 0o1, "USD", "$1.01"},
				{"fr", 1000, 10, "USD", "1\u202f000,10\u00a0$"},
				{"fr-CA", 1000, 10, "USD", "1\u00a0000,10\u00a0$"},
				{"bn", 100000, 1, "USD", "১,০০,০০০.০১$"},
				{"en", 100000, 1, "BHD", "BHD\u00a0100,000.001"},
				{"ar", 100000, 1, "BHD", "\u061c100,000.001\u00a0د.ب.\u200f"},
			} {
				mf := num.MustNewMoneyFormatter(tc.locale)
				_ = mf.SetLocale(tc.locale)
				mf.MustSetLocale(tc.locale)
				mf.UseStandardStyle()
				mf.DisplayCurrencyAsSymbolNarrow()

				actual := mf.MustFormat(tc.whole, tc.frac, tc.cur)
				if actual != tc.expected {
					t.Errorf("test case #%d - got: %v, expected: %v", i+1, actual, tc.expected)
				}
			}
		})

		t.Run("standard style, no currency", func(t *testing.T) {
			for i, tc := range []moneyTestCase{
				{"en", 1, 0o1, "USD", "1.01"},
				{"fr", 1000, 10, "USD", "1\u202f000,10"},
				{"fr-CA", 1000, 10, "USD", "1\u00a0000,10"},
				{"bn", 100000, 1, "USD", "১,০০,০০০.০১"},
				{"en", 100000, 1, "BHD", "100,000.001"},
				{"ar", 100000, 1, "BHD", "\u200f100,000.001"},
			} {
				mf := num.MustNewMoneyFormatter(tc.locale)
				_ = mf.SetLocale(tc.locale)
				mf.MustSetLocale(tc.locale)
				mf.UseStandardStyle()
				mf.DisplayNoCurrency()

				actual := mf.MustFormat(tc.whole, tc.frac, tc.cur)
				if actual != tc.expected {
					t.Errorf("test case #%d - got: %v, expected: %v", i+1, actual, tc.expected)
				}
			}
		})
		t.Run("accounting style, currency code", func(t *testing.T) {
			for i, tc := range []moneyTestCase{
				{"en", -1, 0, "USD", "(USD\u00a01.00)"},
				{"fr", 1000, 10, "USD", "1\u202f000,10\u00a0USD"},
				{"fr", 10000, 10, "USD", "10\u202f000,10\u00a0USD"},
				{"fr", -10000, 10, "USD", "(10\u202f000,10\u00a0USD)"},
				{"ar", -100000, 1, "BHD", "(\u061c100,000.001\u00a0BHD)"},
			} {
				mf := num.MustNewMoneyFormatter(tc.locale)
				_ = mf.SetLocale(tc.locale)
				mf.MustSetLocale(tc.locale)
				mf.UseAccountingStyle()
				mf.DisplayCurrencyAsCode()

				actual := mf.MustFormat(tc.whole, tc.frac, tc.cur)
				if actual != tc.expected {
					t.Errorf("test case #%d - got: %v, expected: %v", i+1, actual, tc.expected)
				}
			}
		})

		t.Run("accounting style, currency symbol", func(t *testing.T) {
			for i, tc := range []moneyTestCase{
				{"en", -1, 0o1, "USD", "($1.01)"},
				{"fr", -1000, 10, "USD", "(1\u202f000,10\u00a0$US)"},
				{"fr-CA", -1000, 10, "USD", "(1\u00a0000,10\u00a0$\u00a0US)"},
				{"en-CA", -1000, 10, "USD", "(US$\u00a01,000.10)"},
				{"ar", -100000, 1, "BHD", "(\u061c100,000.001\u00a0د.ب.\u200f)"},
			} {
				mf := num.MustNewMoneyFormatter(tc.locale)
				_ = mf.SetLocale(tc.locale)
				mf.MustSetLocale(tc.locale)
				mf.UseAccountingStyle()
				mf.DisplayCurrencyAsSymbol()

				actual := mf.MustFormat(tc.whole, tc.frac, tc.cur)
				if actual != tc.expected {
					t.Errorf("test case #%d - got: %v, expected: %v", i+1, actual, tc.expected)
				}
			}
		})

		t.Run("accounting style, currency symbol narrow", func(t *testing.T) {
			for i, tc := range []moneyTestCase{
				{"en", -1, 0o1, "USD", "($1.01)"},
				{"fr", -1000, 10, "USD", "(1\u202f000,10\u00a0$)"},
				{"fr-CA", -1000, 10, "USD", "(1\u00a0000,10\u00a0$)"},
				{"en-CA", -1000, 10, "USD", "($1,000.10)"},
				{"ar", -100000, 1, "BHD", "(\u061c100,000.001\u00a0د.ب.\u200f)"},
			} {
				mf := num.MustNewMoneyFormatter(tc.locale)
				_ = mf.SetLocale(tc.locale)
				mf.MustSetLocale(tc.locale)
				mf.UseAccountingStyle()
				mf.DisplayCurrencyAsSymbolNarrow()

				actual := mf.MustFormat(tc.whole, tc.frac, tc.cur)
				if actual != tc.expected {
					t.Errorf("test case #%d - got: %v, expected: %v", i+1, actual, tc.expected)
				}
			}
		})

		t.Run("accounting style, no currency symbol", func(t *testing.T) {
			for i, tc := range []moneyTestCase{
				{"en", -1, 0o1, "USD", "(1.01)"},
				{"fr", -1000, 10, "USD", "(1\u202f000,10)"},
				{"fr-CA", -1000, 10, "USD", "(1\u00a0000,10)"},
				{"en-CA", -1000, 10, "USD", "(1,000.10)"},
				{"ar", -100000, 1, "BHD", "(100,000.001)"},
			} {
				mf := num.MustNewMoneyFormatter(tc.locale)
				_ = mf.SetLocale(tc.locale)
				mf.MustSetLocale(tc.locale)
				mf.UseAccountingStyle()
				mf.DisplayNoCurrency()

				actual := mf.MustFormat(tc.whole, tc.frac, tc.cur)
				if actual != tc.expected {
					t.Errorf("test case #%d - got: %v, expected: %v", i+1, actual, tc.expected)
				}
			}
		})
	})
}
