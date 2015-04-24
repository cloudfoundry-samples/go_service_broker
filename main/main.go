package main

import (
	"flag"

	conf "github.com/xingzhou/go_service_broker/config"
	util "github.com/xingzhou/go_service_broker/utils"
	webs "github.com/xingzhou/go_service_broker/web_server"
)

func main() {
	// Step0. Get Config Path
	defaultConfigPath := util.GetPath([]string{"assets", "config.json"})
	configPath := flag.String("c", defaultConfigPath, "use '-c' option to specify the config file path")

	// Step1. Load configuration
	_, err := conf.LoadConfig(*configPath)
	if err != nil {
		panic("Error loading config file...")
	}

	// Step2. Load data

	// Step3. Start Server
	server := webs.CreateServer()
	server.Start()
}
