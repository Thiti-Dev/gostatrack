package statrack

import (
	"strconv"
	"strings"
)

//Configs is used to initialize how program should behavior
type configs struct {
	sites    []string
	interval int64
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
