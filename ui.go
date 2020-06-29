package main

import (
	"context"
	"fmt"
	"io"

	"github.com/eiannone/keyboard"
	"github.com/gosuri/uilive"
)

type ui struct {
	wf   writerFlusher
	view *branchesState
}

func newUI(view *branchesState) ui {
	return ui{
		wf:   uilive.New(),
		view: view,
	}
}

func (ui *ui) run(ctx context.Context) error {
	// Setup key input
	keyCh, err := keyboard.GetKeys(10)
	if err != nil {
		return err
	}
	defer keyboard.Close()

	// Listen for keyboard events
	for {
		select {
		case ev := <-keyCh:
			done, err := ui.handleKeyEvent(ev)
			if !done {
				continue
			}
			return err
		case <-ctx.Done():
			return nil
		}
	}
}

func (ui *ui) handleKeyEvent(ev keyboard.KeyEvent) (bool, error) {
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
		ui.view.selectPrevious()
		fmt.Fprintf(ui.wf, "→	%s\n", ui.view.selectedBranchWithColor())
		ui.wf.Flush()

	// Down
	case ev.Key == keyboard.KeyArrowDown:
		fallthrough
	case ev.Key == keyboard.KeyArrowRight:
		ui.view.selectNext()
		fmt.Fprintf(ui.wf, "→	%s\n", ui.view.selectedBranchWithColor())
		ui.wf.Flush()

	// Enter
	case ev.Key == keyboard.KeyEnter:
		return true, cmdRun("git", "checkout", extractBranch(ui.view.selectedBranchName()))

	default:
		return false, nil
	}
	return false, nil
}

type writerFlusher interface {
	io.Writer
	Flush() error
}
