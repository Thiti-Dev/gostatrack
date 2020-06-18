package statrack

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//Configs is used to initialize how program should behavior
type configs struct {
	sites    []string
	interval int64
}

type fetchResponse struct {
	site    string
	isAlive bool
}

//Global declaration
var myConfigs configs
var allSitesInSequence map[string]bool
var totalCountOfFlushes = 0
var isFirstTime = true

func putSitesIntoSequencesMapped() {
	allSitesInSequence = make(map[string]bool)
	for _, site := range myConfigs.sites {
		allSitesInSequence[site] = false
	}
}

//InitializeConfigFromString uses for initializing the configs for futher usages
func InitializeConfigFromString(recieveArgs []string) {
	for _, argStr := range recieveArgs {
		if strings.Contains(argStr, "--site") {
			//Working with site plain args
			siteLists := strings.Split(strings.Split(argStr, "=")[1], ",")
			myConfigs.sites = siteLists // applied to global conf
		} else if strings.Contains(argStr, "--interval") {
			//Working with site plain args
			interval, err := strconv.ParseInt(strings.Split(argStr, "=")[1], 10, 64)
			if err != nil {
				myConfigs.interval = 3000 // default 3 seconds if found an error extracting args
			} else {
				myConfigs.interval = interval
			}
		}
	}
	putSitesIntoSequencesMapped()
}

//LogSettings uses for logging the settings out
func LogSettings() {
	fmt.Println("[", len(myConfigs.sites), " SITES TO TRACK]: ", myConfigs.sites)
	fmt.Printf("[INTERVAL]: %v seconds", myConfigs.interval/1000)
}

func isConnectionAliveFromLink(link string, c chan fetchResponse) {
	_, err := http.Get(link)
	if err != nil {
		c <- fetchResponse{site: link, isAlive: false}
		return
	}
	c <- fetchResponse{site: link, isAlive: true}
}

func generateOutputTextFromStatus(status bool) string {
	var genStr string
	switch status {
	case true:
		genStr = "Server is up [\u001b[32malive\u001b[0m]"
	case false:
		genStr = "Server is down [\u001b[31mnot-alive\u001b[0m]"
	}
	return genStr
}

func outputTrackStatus() {
	totalCountOfFlushes++
	if !isFirstTime {
		for i := -1; i < len(myConfigs.sites); i++ {
			fmt.Printf("\033[A\033[K")
		}
	} else {
		isFirstTime = false
	}
	fmt.Printf("[WAVE]: \u001b[33m%v\u001b[0m\n", totalCountOfFlushes/2)
	for site, status := range allSitesInSequence {
		fmt.Printf("[%v] => status : %v\n", site, generateOutputTextFromStatus(status))
	}
}

func outputCLIHeader() {
	print("\033[H\033[2J")
	fmt.Println("[DEBUG]: Starting tracking . . .")
}

//StartTrackingProcess uses for starting all of the process in tracking a website in interval
func StartTrackingProcess() {
	outputCLIHeader()
	// Creating Channel
	mainChannel := make(chan fetchResponse)

	for _, link := range myConfigs.sites {
		go isConnectionAliveFromLink(link, mainChannel)
	}

	for channelResponse := range mainChannel {
		allSitesInSequence[channelResponse.site] = channelResponse.isAlive // Applied status in memory
		//@Stdout-Phase
		outputTrackStatus()
		// ─────────────────────────────────────────────────────────────────
		go func(c fetchResponse) {
			time.Sleep(time.Duration(myConfigs.interval/1000) * time.Second)
			isConnectionAliveFromLink(c.site, mainChannel)
		}(channelResponse)
	}
}
