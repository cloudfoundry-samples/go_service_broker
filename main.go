package main

import (
	"os"
	"errors"
	"fmt"
	"flag"

	conf "github.com/cloudfoundry-samples/go_service_broker/config"
	utils "github.com/cloudfoundry-samples/go_service_broker/utils"
	webs "github.com/cloudfoundry-samples/go_service_broker/web_server"
)

const (
	AWS = "AWS"
	SOFTLAYER = "SoftLayer"
)

type Options struct {
	ConfigPath string
	Cloud string
}

var options Options

func init() {
	defaultConfigPath := utils.GetPath([]string{"assets", "config.json"})
	flag.StringVar(&options.ConfigPath, "c", defaultConfigPath, "use '-c' option to specify the config file path")

	flag.StringVar(&options.Cloud, "cloud", AWS, "use '--cloud' option to specify the cloud client to use: AWS or SoftLayer")

	flag.Parse()
}

func main() {
	err := checkCloudName(options.Cloud)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	_, err = conf.LoadConfig(options.ConfigPath)
	if err != nil {
		panic(fmt.Sprintf("Error loading config file [%s]...", err.Error()))
	}

	server, err := webs.CreateServer(options.Cloud)
	if err != nil {
		panic(fmt.Sprintf("Error creating server [%s]...", err.Error))
	}

	server.Start()
}

// Private func

func checkCloudName(name string) error {
	fmt.Println(name)

	switch name {
		case AWS, SOFTLAYER:
		return nil
	}

	return errors.New(fmt.Sprintf("Invalid cloud name: %s", name))
}
