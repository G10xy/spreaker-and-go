package main

import (
	"fmt"
	"os"
)


var version = "dev"

func main() {
	if err := Execute(version); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
