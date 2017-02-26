package main

import (
	"fmt"
	"os"

	"github.com/bndw/pick/commands"
	"github.com/bndw/pick/config"
)

// Injected in build
var version string

func main() {
	cfg, err := config.Load(version)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if err := commands.Execute(cfg); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
