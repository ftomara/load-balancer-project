package main

import (
	"bytes"
	"encoding/json"
	"loadbalancer/node"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {
	host_name, err := os.Hostname()
	if err != nil {
		log.Println("error while getting node host name")
		return
	}
	port := os.Getenv("PORT")
	node_url := host_name + port
	weight := rand.Intn(10) + 1
	postData := map[string]any{"url": node_url, "weight": weight, "id": host_name}
	jsonData, _ := json.Marshal(postData)
	for {
		resp, err := http.Post("http://lb:8080/register", "application/json", bytes.NewBuffer(jsonData))
		if err == nil {
			log.Printf("Sccessful Register Request on %s", host_name)
			resp.Body.Close()
			break
		}
		log.Println("LB not ready, retrying in 2 seconds...")
		time.Sleep(2 * time.Second)
	}
	node.StartServer(port)
}
