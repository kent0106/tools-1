# Go Tools

[![PkgGoDev](https://pkg.go.dev/badge/github.com/kent0106/gotools)](https://pkg.go.dev/github.com/kent0106/gotools)

This subrepository holds the source for various packages and tools that support
the Go programming language.

Some of the tools, `godoc` and `vet` for example, are included in binary Go
distributions.

Others, including the Go `guru` and the test coverage tool, can be fetched with
`go get`.

Packages include a type-checker for Go and an implementation of the
Static Single Assignment form (SSA) representation for Go programs.

## Download/Install

The easiest way to install is to run `go get -u github.com/kent0106/gotools/...`. You can
also manually git clone the repository to `$GOPATH/src/github.com/kent0106/gotools`.

## JS/CSS Formatting

This repository uses [prettier](https://prettier.io/) to format JS and CSS files.

The version of `prettier` used is 1.18.2.

It is encouraged that all JS and CSS code be run through this before submitting
a change. However, it is not a strict requirement enforced by CI.

## Report Issues / Send Patches

This repository uses Gerrit for code changes. To learn how to submit changes to
this repository, see https://golang.org/doc/contribute.html.

The main issue tracker for the tools repository is located at
https://github.com/golang/go/issues. Prefix your issue with "x/tools/(your
subdir):" in the subject line, so it is easy to find.
