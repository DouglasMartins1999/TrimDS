package main

import (
	"os"

	"dotins.eu.org/trimds/gui"
	"dotins.eu.org/trimds/lib"
)

func main() {
	if len(os.Args) > 1 {
		lib.AddROMs(os.Args[1:])
		lib.Trim()
	} else {
		gui.Init()
	}
}
