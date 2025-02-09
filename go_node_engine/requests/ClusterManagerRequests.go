package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go_node_engine/logger"
	"go_node_engine/model"
	"io"
	"net/http"
)

/////////////////// GRACEFUL EXIT ///////////////////

// HandshakeAnswer is the struct that describes the handshake answer between the nodes
type ClusterHandshakeAnswer struct {
	MqttPort string `json:"MQTT_BROKER_PORT"`
	NodeId   string `json:"id"`
}

// // Sent to the cluster manager upon exit decision
type ClusterExitRequest struct {
	NodeId     string `json:"node_id"`
	ExitStatus string `json:"exit_status"`
}

type ClusterExitResponse struct {
	Message string `json:"message"`
}

/////////////////// HEURISTICS ///////////////////

// Negotiation Structs
type NegotiationRequest struct {
	NodeId string `json:"node_id"`
}

type NegotiationResponse struct {
	Decisions []JobDecision `json:"decisions"`
}

type JobDecision struct {
	JobID    string `json:"job_id"`
	Decision string `json:"decision"`
}

// Exit Confirmation Structs
type ExitConfirmationRequest struct {
	Message string `json:"message"`
}

type ExitConfirmaionResponse struct {
	Message string `json:"message"`
}

//////////////////////////////////////////////////

// ClusterHandshake sends a handshake request to the cluster manager
func ClusterHandshake(address string, port int) ClusterHandshakeAnswer {
	data, err := json.Marshal(model.GetNodeInfo())
	if err != nil {
		logger.ErrorLogger().Fatalf("Handshake failed, json encoding problem, %v", err)
	}
	jsonbody := bytes.NewBuffer(data)
	resp, err := http.Post(fmt.Sprintf("http://%s:%d/api/node/register", address, port), "application/json", jsonbody)
	if err != nil {
		logger.ErrorLogger().Fatalf("Handshake failed, %v", err)
	}
	if resp.StatusCode != 200 {
		logger.ErrorLogger().Fatalf("Handshake failed with error code %d", resp.StatusCode)
	}
	//defer resp.Body.Close()
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.ErrorLogger().Fatalf("Handshake failed, %v", err)
		}
	}()

	handshakeAnswer := ClusterHandshakeAnswer{}
	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorLogger().Fatalf("Handshake failed, %v", err)
	}
	err = json.Unmarshal(responseBytes, &handshakeAnswer)
	if err != nil {
		logger.ErrorLogger().Fatalf("Handshake failed, %v", err)
	}
	return handshakeAnswer
}

func NotifyClusterExit(address string, port int, node_id string) ClusterExitResponse {
	request := ClusterExitRequest{
		NodeId: node_id,
	}

	data, err := json.Marshal(request)
	if err != nil {
		logger.ErrorLogger().Fatalf("Exit request failed, json encoding problem, %v", err)
	}
	jsonbody := bytes.NewBuffer(data)

	resp, err := http.Post(fmt.Sprintf("http://%s:%d/api/node/request_exit", address, port), "application/json", jsonbody)
	if err != nil {
		logger.ErrorLogger().Fatalf("Exit request failed, %v", err)
	}
	if resp.StatusCode != 200 {
		logger.ErrorLogger().Fatalf("Exit request failed with error code %d", resp.StatusCode)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.ErrorLogger().Fatalf("Exit request failed, %v", err)
		}
	}()

	exitResponse := ClusterExitResponse{}
	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorLogger().Fatalf("Handshake failed, %v", err)
	}
	err = json.Unmarshal(responseBytes, &exitResponse)
	if err != nil {
		logger.ErrorLogger().Fatalf("Handshake failed, %v", err)
	}
	return exitResponse
}

func Negotiate(address string, port int, node_id string) {
	for {

		// Sending Request
		request := NegotiationRequest{
			NodeId: node_id,
		}

		data, err := json.Marshal(request)
		if err != nil {
			logger.ErrorLogger().Fatalf("Exit request failed, json encoding problem, %v", err)
		}
		jsonbody := bytes.NewBuffer(data)

		resp, err := http.Post(fmt.Sprintf("http://%s:%d/api/node/negotiate_exit", address, port), "application/json", jsonbody)
		if err != nil {
			logger.ErrorLogger().Fatalf("Exit request failed, %v", err)
		}
		if resp.StatusCode != 200 {
			logger.ErrorLogger().Fatalf("Exit request failed with error code %d", resp.StatusCode)
		}

		defer func() {
			if err := resp.Body.Close(); err != nil {
				logger.ErrorLogger().Fatalf("Exit request failed, %v", err)
			}
		}()

		// Parsing Response
		negotiationResponse := NegotiationResponse{}

		responseBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.ErrorLogger().Fatalf("Handshake failed, %v", err)
		}

		err = json.Unmarshal(responseBytes, &negotiationResponse)
		if err != nil {
			logger.ErrorLogger().Fatalf("Handshake failed, %v", err)
		}

		fmt.Println("NEGOTIATION RESPONSE:", negotiationResponse.Decisions)
		break
	}

	// return negotiationResponse
}

func ConfirmExit(address string, port int, node_id string) {

}

// func Negotiate(address string, port int, node_id string) {
// 	request := NegotiationRequest{
// 		NodeId: node_id,
// 	}

// 	data, err := json.Marshal(request)
// 	if err != nil {
// 		logger.ErrorLogger().Fatalf("Exit request failed, json encoding problem, %v", err)
// 	}
// 	jsonbody := bytes.NewBuffer(data)

// 	resp, err := http.Post(fmt.Sprintf("http://%s:%d/api/node/request_exit", address, port), "application/json", jsonbody)
// 	if err != nil {
// 		logger.ErrorLogger().Fatalf("Exit request failed, %v", err)
// 	}
// 	if resp.StatusCode != 200 {
// 		logger.ErrorLogger().Fatalf("Exit request failed with error code %d", resp.StatusCode)
// 	}

// 	defer func() {
// 		if err := resp.Body.Close(); err != nil {
// 			logger.ErrorLogger().Fatalf("Exit request failed, %v", err)
// 		}
// 	}()

// 	exitResponse := ExitConfirmaionResponse{}
// 	responseBytes, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		logger.ErrorLogger().Fatalf("Handshake failed, %v", err)
// 	}
// 	err = json.Unmarshal(responseBytes, &exitResponse)
// 	if err != nil {
// 		logger.ErrorLogger().Fatalf("Handshake failed, %v", err)
// 	}
// }
