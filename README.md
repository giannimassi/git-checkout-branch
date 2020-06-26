![Build](https://github.com/giannimassi/git-checkout-branch/workflows/Test/badge.svg?branch=master)

# git-checkout-branch

## A cli application for interactively checking out branches

[![asciicast](https://asciinema.org/a/6nK7EchQbDNJzD9y0Cb2dIPgv.svg)](https://asciinema.org/a/6nK7EchQbDNJzD9y0Cb2dIPgv)

## Installation

```sh
go get -u github.com/giannimassi/git-checkout-interactive
```

If $GOBIN is available in your $PATH the command can be used as `git checkout-interactive`.

### Optional

You can also add the following to your `.gitconfig` to provide a simpler to use alias for this command:

```ini
[alias]
    cb = checkout-interactive --sort=-committerdate
```

After adding this alias you can simply call this to choose a branch from the latest ones:

```sh
git cb
```

This is just an example, you can add flags and arguments as needed in the alias.

## Usage

```sh
git checkout-interactive <git branch arguments and flags>
```

Displays the current branch and allows to circle through branches listed with `git branch` and the provided flag and arguments.

- press down/left/j to select next branch
- press up/right/h to select previous branch
- press enter to checkout selected branch
- press ESC or ctrl-c to exit without any changes

### Examples

```sh
git checkout-interactive
```

Choose a branch to checkout among the local branches.

```sh
git checkout-interactive -a
```

Choose a branch to checkout among all branches, both local and remote.

```sh
git checkout-interactive --sort=-committerdate
```

Choose a branch to checkout among the local branches listed by most recent commit.


#### Have fun versioning!