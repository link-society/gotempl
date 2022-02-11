package main

import (
	"fmt"
	"os"

	"github.com/link-society/gotempl/internal/io"
)

func main() {
	err := io.ExecuteTemplate(os.Args[1:])
	if err != nil {
		fmt.Fprintln(io.Stderr(), err)
		os.Exit(1)
	}
}
