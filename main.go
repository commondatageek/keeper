package main

import (
	"fmt"
	"github.com/commondatageek/keeper/lib"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintln(os.Stderr, "Usage: keeper {filename} {test-data} {number}")
		os.Exit(1)
	}

	var filename string = os.Args[1]
	var url string = os.Args[2]
	var n string = os.Args[3]

	intN, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		log.Fatalf("Could not convert %s to integer", n)
	}

	// construct ItemList
	items := make(lib.ItemList, intN)
	var i int64
	for i = 0; i < intN; i++ {
		var w lib.WebSite = lib.NewWebSite(url)
		items[i] = w
	}

	db := lib.NewLocalDatabase(filename)
	writeN, writeErr := db.Write(items)
	if writeErr != nil {
		log.Fatalf("Error writing database file: %s", writeErr)
	}

	log.Printf("Wrote %d items\n", writeN)
}
