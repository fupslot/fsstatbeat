package main

import (
	"os"

	"github.com/fupslot/fsstatbeat/cmd"

	_ "github.com/fupslot/fsstatbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
