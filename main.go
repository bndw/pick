package main

import (
	"fmt"
	"os"

	"github.com/bndw/pick/commands"
)

const Version = "v0.2.2"

func main() {
	if err := commands.Execute(Version); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
