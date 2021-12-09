package main

import (
	"fmt"
	"os"

	"github.com/link-society/gotempl/internal/io"
	"github.com/link-society/gotempl/internal/options"
)

func ExecuteTemplate(args []string) error {
	opts, err := options.ParseOptions(args)
	if err != nil {
		return err
	}

	context, err := io.NewContext(opts)
	if err != nil {
		return err
	}

	err = context.Write()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := ExecuteTemplate(os.Args[1:])
	if err != nil {
		fmt.Fprintln(io.Stderr(), err)
		os.Exit(1)
	}
}
