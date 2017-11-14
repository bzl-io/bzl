package bzl

import (
	"testing"
	"github.com/bzl-io/bzl"
)

func TestFix(t *testing.T) {
	app := bzl.New()
	app.Run([]string {"help"})
	// Output: foo
}
