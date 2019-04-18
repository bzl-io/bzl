package main

import (
	"os"
)

func main() {
	app := NewApp()
	app.Run(os.Args)
}
