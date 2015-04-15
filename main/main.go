package main

import (
	"fmt"

	aws "github.com/awslabs/aws-sdk-go/aws"
	ac "github.com/xingzhou/go_service_broker/aws_client"
	. "github.com/xingzhou/go_service_broker/web_server"
)

func main() {
	server := CreateServer()

	client := ac.NewClient("us-east-1")

	// output, err := client.CreateInstance()
	// handleAWSError("CreateInstance", output, err)

	output, err := client.CreateKeyPair("mykey1")
	handleAWSError("CreateKeyPair", output, err)

	server.Start()
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
