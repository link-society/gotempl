package io

import (
	"io"
	"os"
)

type StdFiles struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

var std = NewStdFiles()

func NewStdFiles() *StdFiles {
	return &StdFiles{
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}

func SetInput(stdin io.Reader) {
	std.stdin = stdin
}

func SetOutput(stdout io.Writer) {
	std.stdout = stdout
}

func SetErrOutput(stderr io.Writer) {
	std.stderr = stderr
}

func Stdin() io.Reader {
	return std.stdin
}

func Stdout() io.Writer {
	return std.stdout
}

func Stderr() io.Writer {
	return std.stderr
}
