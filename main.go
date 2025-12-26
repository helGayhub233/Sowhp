package main

import (
	"Sowhp/core"
	"os"
)

func main() {
	if err := core.Run(); err != nil {
		os.Exit(1)
	}
}
