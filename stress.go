package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	authURL         = "http://localhost:8081/login"
	simulationURL   = "http://localhost:8081/post-simulation"
	username        = "admin"
	password        = "password"
	concurrentUsers = 223
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	AuthDataResponse AuthDataResponse `json:"data"`
}
type AuthDataResponse struct {
	Token string `json:"token"`
}

type SimulationRequest struct {
	ScenarioID                int    `json:"scenario_id"`
	InputType                 string `json:"input_type"`
	ScenarioType              string `json:"scenario_type"`
	Role                      string `json:"role"`
	SimulationAdvancementType string `json:"simulation_advancement_type"`
	Mode                      string `json:"mode"`
}

func authenticate() (string, error) {
	reqBody, _ := json.Marshal(AuthRequest{Username: username, Password: password})
	resp, err := http.Post(authURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("auth error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("auth failed: %s", string(body))
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return "", fmt.Errorf("invalid auth response: %v", err)
	}

	return authResp.AuthDataResponse.Token, nil
}

func simulate(token string, id int, wg *sync.WaitGroup) {
	defer wg.Done()

	payload := SimulationRequest{
		ScenarioID:                1,
		InputType:                 "text",
		ScenarioType:              "takeoff",
		Role:                      "tower",
		SimulationAdvancementType: "steps",
		Mode:                      "single",
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", simulationURL, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("[User %d] request creation error: %v\n", id, err)
		return
	}
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[User %d] request failed: %v\n", id, err)
		return
	}
	defer resp.Body.Close()

	//log.Printf("[User %d] Response: %s\n", id, respBody)
}

func main() {
	token, err := authenticate()
	if err != nil {
		log.Fatalf("Failed to authenticate: %v\n", err)
	}
	log.Println("Authenticated successfully. Starting stress test...")

	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go simulate(token, i+1, &wg)
	}

	wg.Wait()
	duration := time.Since(start)
	log.Printf("Stress test completed in %s\n", duration)
}
