// Package version allows you to fetch the CLDR version used in data generation.
package version

import "github.com/ttzhou/cldr/internal/locale"

// Get returns the CLDR revision used by this module for its various packages.
// This can update at any point as the package evolves.
func Get() string {
	return locale.CLDRVersion
}
