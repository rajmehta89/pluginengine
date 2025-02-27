package server

import (
	"NMS/src/plugin/windows"
	"encoding/json"
	"fmt"
	"github.com/pebbe/zmq4"
	"log"
)

type Request struct {
	RequestType string `json:"RequestType"`
	SystemType  string `json:"SystemType"`
	Ip          string `json:"Ip"`
	Username    string `json:"Username"`
	Password    string `json:"password"`
}

func handleRequest(requestStr string) string {

	fmt.Println(requestStr)
	var req Request

	err := json.Unmarshal([]byte(requestStr), &req)
	fmt.Println("after conv", req)

	if err != nil {

		log.Println("Invalid JSON format received")

		return `{"error": "Invalid JSON format"}`
	}

	switch req.RequestType {

	case "discovery":
		return handleDiscovery(req)

	case "provisioning":
		return handleProvisioning(req)

	default:

		log.Println("Received unknown request type")

		return `{"error": "Unknown request type"}`

	}
}

func handleDiscovery(req Request) string {

	var metrics interface{}

	switch req.SystemType {

	case "windows":

		fmt.Println("calling discover with "+req.Ip, req.Username, req.Password)
		metrics = windows.Discover(req.Ip, req.Username, req.Password)

	default:
		return `{"error": "Unknown discovery type"}`
	}

	res, _ := json.Marshal(metrics)

	log.Println("Discovery completed for:", req.SystemType)

	return string(res)

}

func handleProvisioning(req Request) string {

	var response interface{}

	switch req.SystemType {

	case "windows":
		response = windows.Start(req.Ip, req.Username, req.Password)

	default:

		return `{"error": "Unknown provisioning type"}`

	}

	res, _ := json.Marshal(response)

	log.Println("Provisioning completed for:", req.SystemType)

	return string(res)

}

func sendResponse(socket *zmq4.Socket, response string) {

	_, err := socket.Send(response, 0)

	if err != nil {

		log.Println("Error sending response:", err)
	}

}

func StartZMQServer() {

	socket, _ := zmq4.NewSocket(zmq4.REP)

	defer socket.Close()

	socket.Bind("tcp://*:5555")

	for {

		requestStr, err := socket.Recv(0)

		if err != nil {

			log.Println("Error receiving request:", err)

			continue
		}

		log.Println("Received request:", requestStr)

		response := handleRequest(requestStr)

		sendResponse(socket, response)

	}
}
