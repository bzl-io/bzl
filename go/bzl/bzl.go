package main

import (
	"os"
	"github.com/bzl-io/bzl"
)

func main() {
	app := bzl.New()
	app.Run(os.Args)
}
