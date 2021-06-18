package main

import (
	"encoding/json"
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
	var data string = os.Args[2]

	jsonBytes, marshalErr := json.Marshal(data)
	if marshalErr != nil {
		log.Fatalf("Could not marshal JSON: %s\n", marshalErr)
	}

	err := lib.SafeWriteFile(filename, &jsonBytes)
	if err != nil {
		log.Fatalf("Could not write file %s: %s\n", filename, err)
	}
}
