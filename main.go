package main

import (
	"os"

	"github.com/juhc/contos/cmd/control"
	"github.com/juhc/contos/cmd/initcontainers"
)

var entrypoints = map[string]func(){
	"autologin":          control.AutologinMain,
}

func main() {

	initcontainers.Init()

	if len(os.Args) > 1 {
		if f, ok := entrypoints[os.Args[1]]; ok {
			f()
			return
		}
	}

	control.Main()
}
