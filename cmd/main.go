package main

import (
	"fmt"
	"io"
	"os"

	"github.com/adityanatraj/dronevis"
)

func main() {
	args := os.Args[1:]

	var rdr io.Reader
	if len(args) == 0 {
		rdr = os.Stdin
	} else if len(args) == 1 {
		f, err := os.Open(args[0])
		if err != nil {
			errorOut(err)
		} else {
			defer f.Close()
		}
		rdr = f
	} else {
		errorOut(fmt.Errorf("multiple files specified. sorry. %s", "ugh"))
	}

	output, err := dronevis.Graph2(rdr)
	if err != nil {
		errorOut(err)
	}
	fmt.Println(output[1])
}

func errorOut(err error) {
	fmt.Printf("[error]: %s\n", err)
	os.Exit(1)
}
