package main

import (
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

func cmdOutput(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	b, err := cmd.CombinedOutput()
	return string(b), err
}

func cmdRun(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// branch describes a Git branch.
type branch struct {
	name      string
	isCurrent bool
}

// branchesState provides the state for the application, based on CQRS
// Includes view and command methods
type branchesState struct {
	branches []branch
	selected int
}

func (b *branchesState) init(args []string) error {
	args = append([]string{"branch"}, args...)
	out, err := cmdOutput("git", args...)
	if err != nil {
		println(out)
		return err
	}
	b.branches = splitBranches(out)
	b.selectCurrent()
	return nil
}

func splitBranches(output string) []branch {
	names := strings.Split(output, "\n")
	var branches []branch
	for _, name := range names {
		if len(name) == 0 {
			continue
		}
		isCurrent := false
		if strings.Contains(name, "*") {
			name = strings.Replace(name, "*", "", -1)
			isCurrent = true
		}

		name = strings.TrimSpace(name)
		branches = append(branches, branch{name: name, isCurrent: isCurrent})
	}
	return branches
}

func extractBranch(name string) string {
	if strings.Contains(name, "->") {
		s := strings.Split(name, "->")
		return strings.TrimSpace(s[0])
	}
	return name
}

// Commands
// These methods allow to mutate state

func (b *branchesState) selectCurrent() {
	for ; !b.branches[b.selected].isCurrent; b.selected++ {
	}
}

func (b *branchesState) selectNext() {
	b.selected = (b.selected + len(b.branches) - 1) % len(b.branches)
}

func (b *branchesState) selectPrevious() {
	b.selected = ((b.selected + 1) % len(b.branches))
}

// View
// These methods allow to present the state

func (b *branchesState) selectedBranchName() string {
	return b.branches[b.selected].name
}

func (b *branchesState) selectedBranchWithColor() string {
	formatted := []string{}
	fields := strings.Split(b.selectedBranchName(), "/")
	for _, f := range fields {
		formatted = append(formatted, withColor(f, colorDefaults))
	}
	out := strings.Join(formatted, withColor("/", colorDefaults))
	if b.branches[b.selected].isCurrent {
		out += "\t(" + withColor("currently checked-out", colorDefaults) + ")"
	}
	return out
}

type formatFunc func(format string, a ...interface{}) string

type colorFormatter struct {
	pattern *regexp.Regexp
	format  formatFunc
}

func newColorFormatter(pattern string, format formatFunc) colorFormatter {
	return colorFormatter{
		pattern: regexp.MustCompile(pattern),
		format:  format,
	}
}

func withColor(str string, colors []colorFormatter) string {
	for _, f := range colors {
		if f.pattern.MatchString(str) {
			return f.format(str)
		}
	}
	return str
}

var colorDefaults = []colorFormatter{
	newColorFormatter("feature", color.GreenString),
	newColorFormatter("test", color.GreenString),

	newColorFormatter("improvement", color.YellowString),
	newColorFormatter("refactor", color.YellowString),
	newColorFormatter("currently checked-out", color.YellowString),

	newColorFormatter("fix", color.RedString),
	newColorFormatter("bugfix", color.RedString),
	newColorFormatter("bug", color.RedString),
	newColorFormatter("fixup", color.RedString),
	newColorFormatter("debug", color.RedString),
	newColorFormatter("backup", color.RedString),

	newColorFormatter("master", color.CyanString),

	newColorFormatter("remotes", color.MagentaString),
	newColorFormatter("origin", color.MagentaString),

	newColorFormatter(".*", color.WhiteString),
}
