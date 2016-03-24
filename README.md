# git-extern-sync

**Warning: this is alpha software and should not be used yet!**

This little `go` package lets you keep "external" files (for now, HTTP
accessible) synchronized without using `git-subtree` or `git-submodules`.

## Installation

```
$ go get -u github.com/tjk/git-extern-sync
```

This will download the source into
`$GOPATH/src/github.com/tjk/git-extern-sync` and build and install a
binary at `$GOPATH/bin/git-extern-sync`.

## Usage

`git-extern-sync` leverages `.gitignore` files in the repository in
which it is called to understand which external files you'd like to
syncrhonize.

Simply add a "metadata" line above the target file in `.gitignore` of
the form `sync:<url>` so that it is synchronized.

Given the following `.gitignore`:

```
*.db
# sync:https://raw.githubusercontent.com/tjk/home/2f6fb49bc025cf4907377d99df39abe593ca5890/README.md
file-to-keep-synced
```

Calling `git-extern-sync` alongside it will ensure the untracked file
at `./file-to-keep-synced` in the repo matches the one at the URL.

## Why?

For one, I found this to be a nice way to manage shared `.proto`
definitions.

## TODO

- [ ] support shortcut, ie. `github.com/<name>/<repo>@<ref>` ?
- [ ] quiet flag so no ./gitignore found is silent?
- [ ] add tests
