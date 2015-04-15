package main

import (
	. "github.com/xingzhou/go_service_broker/web_server"
)

func main() {
	server := CreateServer()
	server.Start()
}
