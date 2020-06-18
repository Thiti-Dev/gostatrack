package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"./helpers/statrack"
)

func main() {
	numOfSites := len(os.Args) - 1
	if numOfSites <= 0 {
		fmt.Println("[USAGE]: => go run main.go --site=http://google.com,http://twitter.com --interval=5000")
	}
	statrack.InitializeConfigFromString(os.Args[1:])
	statrack.LogSettings()
	fmt.Print("Proceed tracking? [TYPE y OR n] : ")
	reader := bufio.NewReader(os.Stdin) //create new reader
	actionString, _ := reader.ReadString('\n')
	if strings.Contains(actionString, "y") {
		// Proceed the tracker
		statrack.StartTrackingProcess()
	} else {
		fmt.Println("Exiting .... ")
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}
}
