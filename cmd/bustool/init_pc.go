//go:build !tinygo
// +build !tinygo

package main

import (
	"os"
)

func configIO() (io.Reader, io.Writer) { return os.Stdin, os.Stdout }
