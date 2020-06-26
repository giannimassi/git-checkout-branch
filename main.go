package main

import (
	"context"
	"os"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		println("Unexpected error: " + err.Error())
	}
}

func run(args []string) error {
	// Get branches
	view := &branchesState{}
	if err := view.init(args); err != nil {
		return err
	}

	// Run ui
	ui := newUI(view)
	return ui.run(context.Background())
}
