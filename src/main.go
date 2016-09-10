package main

import (
	"fmt"
	"os"

	"github.com/bndw/pick/commands"
	"github.com/bndw/pick/config"
)

const Version = "v0.2.2"

func main() {
	cfg, err := config.Load(Version)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if err := commands.Execute(cfg); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
