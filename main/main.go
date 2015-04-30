package main

import (
	"flag"

	conf "github.com/xingzhou/go_service_broker/config"
	utils "github.com/xingzhou/go_service_broker/utils"
	webs "github.com/xingzhou/go_service_broker/web_server"
)

func main() {
	// Step1. Get Config Path
	defaultConfigPath := utils.GetPath([]string{"assets", "config.json"})
	configPath := flag.String("c", defaultConfigPath, "use '-c' option to specify the config file path")

	// Step2. Load configuration
	_, err := conf.LoadConfig(*configPath)
	if err != nil {
		panic("Error loading config file...")
	}

	// Step3. Start Server
	server := webs.CreateServer()
	server.Start()
}
