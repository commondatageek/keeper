package main

import (
	"bytes"
	"fmt"
	"github.com/commondatageek/keeper/lib"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: keeper {filename} {test-data}")
		os.Exit(1)
	}

	var filename string = os.Args[1]
	var url string = os.Args[2]

	blah := []byte(url)
	reader := bytes.NewReader(blah)

	writeN, writeErr := lib.SafeWriteFile(filename, reader)
	if writeErr != nil {
		log.Fatalf("Error writing database file: %s", writeErr)
	}

	log.Printf("Wrote %d bytes\n", writeN)
}
