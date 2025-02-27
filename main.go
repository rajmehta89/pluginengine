package main

import (
	"NMS/src/server"
	"NMS/src/util"
)

func main() {

	logger := util.InitializeLogger()

	logger.LogInfo("Starting ZeroMQ Server...")

	server.StartZMQServer()

	logger.LogInfo("ZeroMQ Server Stopped.")

}
