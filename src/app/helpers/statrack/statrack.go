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
		genStr = "Server is up [alive]"
	case false:
		genStr = "Server is down [not-alive]"
	}
	return genStr
}

//StartTrackingProcess uses for starting all of the process in tracking a website in interval
func StartTrackingProcess() {
	fmt.Println("[DEBUG]: Starting tracking . . .")

	// Creating Channel
	mainChannel := make(chan fetchResponse)

	for _, link := range myConfigs.sites {
		go isConnectionAliveFromLink(link, mainChannel)
	}

	for channelResponse := range mainChannel {
		fmt.Println("[", channelResponse.site, "] => status: ", generateOutputTextFromStatus(channelResponse.isAlive))
		go func(c fetchResponse) {
			time.Sleep(time.Duration(myConfigs.interval/1000) * time.Second)
			isConnectionAliveFromLink(c.site, mainChannel)
		}(channelResponse)
	}
}
