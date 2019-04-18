package main

import (
	"testing"
)

func TestFix(t *testing.T) {
	app := NewApp()
	app.Run([]string{"help"})
	// Output: foo
}
