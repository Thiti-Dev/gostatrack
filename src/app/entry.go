package main

import (
	"fmt"
	"os"
)

func main() {
	numOfSites := len(os.Args) - 1
	if numOfSites <= 0 {
		fmt.Println("[USAGE]: => go run main.go --site=http://google.com,http://twitter.com --interval=5000")
	}
}
