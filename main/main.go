package main

import (
	"flag"
	"fmt"
	"time"

	aws "github.com/awslabs/aws-sdk-go/aws"
	ac "github.com/xingzhou/go_service_broker/aws_client"
	conf "github.com/xingzhou/go_service_broker/config"
	util "github.com/xingzhou/go_service_broker/utils"
	webs "github.com/xingzhou/go_service_broker/web_server"
)

func main() {
	// Step1. Get Config Path
	defaultConfigPath := util.GetPath([]string{"assets", "config.json"})
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

func test_aws() {
	client := ac.NewClient("us-east-1")

	// instanceId, err := client.CreateInstance()
	// handleAWSError("CreateInstance", instanceId, err)
	instanceId := "i-fa6a9c2c"

	// handleAWSError("CreateKeyPair", publicKey, err)
	for {
		state, err := client.GetInstanceState(instanceId)
		handleAWSError("GetInstanceStatus", state, err)
		if state == "running" {
			break
		}
		time.Sleep(time.Duration(1) * time.Second)
	}

	privateKey, err := client.InjectKeyPair(instanceId)
	handleAWSError("InjectKeyPair", privateKey, err)

	err = client.RevokeKeyPair(instanceId, privateKey)
	handleAWSError("RevokeKeyPair", "", err)

	//err = client.DeleteInstance(instanceId)
	//handleAWSError("DeleteInstance", "", err)

	// output, err := client.CreateKeyPair("mykey1")
	// handleAWSError("CreateKeyPair", output, err)

	//server.Start()

	// kd, err := utils.ReadFile("/Users/huazhang/.gsb/broker_id_rsa")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// publicKey, err := utils.GeneratePublicKey(kd)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(publicKey)

}

func handleAWSError(operation string, output string, err error) {
	if awserr := aws.Error(err); awserr != nil {
		// A service error occurred.
		fmt.Println("Error:", awserr.Code, awserr.Message)
	} else if err != nil {
		// A non-service error occurred.
		panic(err)
	}
	fmt.Sprintln("Output of %s:", operation)
	fmt.Println(output)
}
