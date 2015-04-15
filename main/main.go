package main

import (
	// "fmt"

	// ac "github.com/xingzhou/go_service_broker/aws_client"
	. "github.com/xingzhou/go_service_broker/web_server"
)

func main() {
	server := CreateServer()

	// client := ac.NewClient("us-east-1")
	// err := client.CreateInstance()
	// if err != nil {
	// 	fmt.Println(err)
	// 	panic(err)
	// }

	server.Start()
}
