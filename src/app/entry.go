package main

import (
	"fmt"
	"os"

	"./helpers/statrack"
)

func main() {
	numOfSites := len(os.Args) - 1
	if numOfSites <= 0 {
		fmt.Println("[USAGE]: => go run main.go --site=http://google.com,http://twitter.com --interval=5000")
	}
	statrack.InitializeConfigFromString(os.Args[1:])
}
