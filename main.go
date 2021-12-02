package main

import (
	"fmt"
	"os"

	"github.com/link-society/gotempl/internal"
)

func main() {
	opts, err := internal.NewOptions(os.Args[1:])

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	context, err := opts.ReadInputFiles()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = opts.WriteOutput(context)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
