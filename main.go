package main

import (
	"fmt"
	"github.com/commondatageek/keeper/lib"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: keeper {filename} {url}")
		os.Exit(1)
	}

	var filename string = os.Args[1]
	var url string = os.Args[2]

	// get list of current items
	db := lib.NewLocalDatabase(filename)
	items, readError := db.Read()
	if readError != nil {
		log.Fatalf("Could not read database: %s\n", readError)
	}

	// create and add our new item
	newItem := lib.NewWebSite(url)
	items = append(items, newItem)

	// write items out to file
	writeN, writeErr := db.Write(items)
	if writeErr != nil {
		log.Fatalf("Error writing database file: %s", writeErr)
	}

	log.Printf("Wrote %d items\n", writeN)
}
