package main

import (
	"os"

	"github.com/nelsonjchen/google-cloud-nuke/v1/cmd"
)

func main() {
	if err := cmd.NewRootCommand().Execute(); err != nil {
		os.Exit(-1)
	}
}
