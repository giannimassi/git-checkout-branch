package main

import (
	"fmt"
	"io"
	"os"

	"github.com/eiannone/keyboard"
	"github.com/gosuri/uilive"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		println("Unexpected error: " + err.Error())
	}
}

func run(args []string) error {
	// Init state
	view := &branchesState{}
	if err := view.init(args); err != nil {
		return err
	}

	// Setup key input
	keyCh, err := keyboard.GetKeys(10)
	if err != nil {
		return err
	}
	defer keyboard.Close()

	// Setup writer and start with currently selected branch
	stdout := uilive.New()
	fmt.Fprintf(stdout, "→	%s\n", view.selectedBranchWithColor())
	stdout.Flush()

	// Listen for keyboard events and term signal
	for {
		select {
		case ev := <-keyCh:
			done, err := handleKeyEvent(stdout, view, ev)
			if !done {
				continue
			}
			return err
		}
	}
}

func handleKeyEvent(wf writerFlusher, view *branchesState, ev keyboard.KeyEvent) (bool, error) {
	if ev.Err != nil {
		return true, ev.Err
	}
	switch {
	// Exit
	case ev.Key == keyboard.KeyCtrlC:
		fallthrough
	case ev.Key == keyboard.KeyEsc:
		return true, nil

	// Up
	case ev.Key == keyboard.KeyArrowUp:
		fallthrough
	case ev.Key == keyboard.KeyArrowLeft:
		fallthrough
	case ev.Rune == 'h':
		view.selectPrevious()
		fmt.Fprintf(wf, "→	%s\n", view.selectedBranchWithColor())
		wf.Flush()

	// Down
	case ev.Key == keyboard.KeyArrowDown:
		fallthrough
	case ev.Key == keyboard.KeyArrowRight:
		fallthrough
	case ev.Rune == 'j':
		view.selectNext()
		fmt.Fprintf(wf, "→	%s\n", view.selectedBranchWithColor())
		wf.Flush()

	// Enter
	case ev.Key == keyboard.KeyEnter:
		return true, cmdRun("git", "checkout", extractBranch(view.selectedBranchName()))

	default:
		return false, nil
	}
	return false, nil
}

type writerWrapper struct {
	io.Writer
}

func (w *writerWrapper) Flush() error {
	return nil
}

type writerFlusher interface {
	io.Writer
	Flush() error
}
