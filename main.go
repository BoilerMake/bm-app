package main

import (
	"os"

	"github.com/cjenright/new-backend/cmd"
)

func main() {
	// NOTE this could be done with something more complex like cobra or another CLI framework
	// But that's a lot of overhead for something that can be pretty straightforward
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "serve" {
		cmd.Serve()
	}
}
