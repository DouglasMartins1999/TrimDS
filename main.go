package main

import (
	_ "embed"
	"os"

	_ "dotins.eu.org/trimds/syso"

	"dotins.eu.org/trimds/gui"
	"dotins.eu.org/trimds/lib"
)

//go:embed assets/icon.ico
var icon []byte

func main() {
	if len(os.Args) > 1 {
		lib.AddROMs(os.Args[1:])
		lib.Trim()
	} else {
		gui.Init(icon)
	}
}
