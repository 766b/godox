GoDoX
===

[![Build Status](https://travis-ci.org/matoous/go-nanoid.svg?branch=master)](https://travis-ci.org/matoous/godox)
[![GoDoc](https://godoc.org/github.com/matoous/godox?status.svg)](https://godoc.org/github.com/matoous/godox)
[![Go Report Card](https://goreportcard.com/badge/github.com/matoous/godox)](https://goreportcard.com/report/github.com/matoous/godox)
[![GitHub issues](https://img.shields.io/github/issues/matoous/godox.svg)](https://github.com/matoous/godox/issues)
[![License](https://img.shields.io/badge/license-MIT%20License-blue.svg)](https://github.com/matoous/godox/LICENSE)

GoDoX extracts comments from Go code based on keywords.

Installation
---

    go get -u github.com/matoous/godox

Usage
---
Any comment lines starting with `TODO` or `FIX` (or other specified keywords, case insensitive) are extracted. 
If `TODO`/`FIX` comments is longer that 1 line, then only first line will be extracted.

    $ godox [<flags>] [<path>...]

    $ godox -h
        Usage: godox [path] [path] ...
          -h    Print help
          -keys string
                Change keywords (default "todo,bug,fix")

    $ godox ./path/to/directory ./path/to/secondary/directory
    example.go:3:1:TODO: Implement io.Writer interface
    example.go:7:14:TODO: Rename field
    example.go:10:1:TODO: Add JWT verification
    example.go:15:3:FIX: Something that is broken
