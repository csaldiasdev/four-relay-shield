package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"syscall"
	"time"

	"github.com/mitchellh/mapstructure"
)

func main() {

	var uname syscall.Utsname
	if err := syscall.Uname(&uname); err != nil {
		fmt.Printf("Uname: %v", err)
	}

	println(cmdTitle)
	print("Listen in port	: ")
	println("8080")

	print("OS		: ")
	println(arrayToString(uname.Sysname))

	print("Kernel		: ")
	println(arrayToString(uname.Release))

	print("Arch		: ")
	println(arrayToString(uname.Machine))

	println("\n███████\n")

	InitController()
	InitWebsocketServer()

	go notifyRelayStatus()

	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/ws", websocketRequest)
	http.ListenAndServe(":8080", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	bodyContent := []byte(indexHTML)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(bodyContent)
}

func websocketRequest(w http.ResponseWriter, r *http.Request) {
	wsServer := GetWebsocketServerInstance()
	idConnection, _ := wsServer.AddConnection(w, r, wsHandleRequest)

	controller := GetControllerInstance()

	relay1, relay2, relay3, relay4 := controller.GetAllRelaysStatus()

	status := RelaysStatus{
		Relay1Status: relay1,
		Relay2Status: relay2,
		Relay3Status: relay3,
		Relay4Status: relay4,
	}

	bodyBytes := prepareMessageToSocket("INITIAL_RELAY_STATUS", status)

	print("Server message ==> ")
	println(string(bodyBytes))

	wsServer.SendMessage(bodyBytes, idConnection)
}

func wsHandleRequest(body []byte) {

	print("Client message ==> ")
	println(string(body))

	var clientMessage []interface{}

	dec := json.NewDecoder(bytes.NewReader(body))

	err := dec.Decode(&clientMessage)

	if err != nil {
		return
	}

	messageTopic, ok := clientMessage[0].(string)

	if !ok {
		return
	}

	switch messageTopic {

	case "SET_VALUE_RELAY_COMMAND":
		var setValueRelayCommand SetValueRelayCommand

		mapstructure.Decode(clientMessage[1], &setValueRelayCommand)

		relayInstance := GetControllerInstance()
		switch setValueRelayCommand.RelayNumber {

		case 1:
			relayInstance.SetRelay1Status(setValueRelayCommand.RelayValue)
			break
		case 2:
			relayInstance.SetRelay2Status(setValueRelayCommand.RelayValue)
			break
		case 3:
			relayInstance.SetRelay3Status(setValueRelayCommand.RelayValue)
			break
		case 4:
			relayInstance.SetRelay4Status(setValueRelayCommand.RelayValue)
			break
		default:
			return
		}

		break
	default:
		return
	}

}

func notifyRelayStatus() {
	controller := GetControllerInstance()
	wsServer := GetWebsocketServerInstance()

	for status := range controller.changeStatusChannel {

		messageStatus := RelayStatusChange{
			RelayNumber:   status.relayNumber,
			RelayNewValue: status.relayNewValue,
			ChangedAt:     status.changedAt,
		}

		bodyBytes := prepareMessageToSocket("NOTIFY_RELAY_STATUS", messageStatus)

		print("Server message ==> ")
		println(string(bodyBytes))

		wsServer.BroadcastMessage(bodyBytes)
	}
}

func prepareMessageToSocket(topic string, data interface{}) []byte {
	message := make([]interface{}, 0, 2)
	message = append(message, topic, data)

	jsonBytes, _ := json.Marshal(message)

	return jsonBytes
}

// RelaysStatus struct
type RelaysStatus struct {
	Relay1Status bool `json:"relay1Status"`
	Relay2Status bool `json:"relay2Status"`
	Relay3Status bool `json:"relay3Status"`
	Relay4Status bool `json:"relay4Status"`
}

// RelayStatusChange struct
type RelayStatusChange struct {
	RelayNumber   int       `json:"relayNumber"`
	RelayNewValue bool      `json:"relayNewValue"`
	ChangedAt     time.Time `json:"changedAt"`
}

// SetValueRelayCommand struct
type SetValueRelayCommand struct {
	RelayNumber int  `json:"relayNumber"`
	RelayValue  bool `json:"relayValue"`
}

func arrayToString(x [65]int8) string {
	var buf [65]byte
	for i, b := range x {
		buf[i] = byte(b)
	}
	str := string(buf[:])
	if i := strings.Index(str, "\x00"); i != -1 {
		str = str[:i]
	}
	return str
}
