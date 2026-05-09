package version_test

import (
	"testing"

	"github.com/ttzhou/cldr/internal/locale"
	"github.com/ttzhou/cldr/version"
)

func Test(t *testing.T) {
	t.Run("Get()", func(t *testing.T) {
		usedForGen, usedForModule := locale.CLDRVersion, version.Get()

		if usedForGen != usedForModule {
			t.Errorf(
				"module CLDR version %v does not match version used to generate CLDR data %v",
				usedForModule,
				usedForGen,
			)
		}
	})
}
