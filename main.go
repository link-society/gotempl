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

	context, err := internal.ReadInputFiles(opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = internal.WriteOutput(opts, context)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
