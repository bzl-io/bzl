package main

import (
	"github.com/bzl-io/bzl"
	"os"
)

func main() {
	app := bzl.New()
	app.Run(os.Args)
}
