package lb

import (
	"log"
	"net/http"
	"time"
)

func StartHealthChecking() {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		checkHealth()
	}
}

func checkHealth() {
	client := http.Client{Timeout: 2 * time.Second}
	for i := range nodes_map {
		modify_mu.Lock()
		url := nodes_map[i].Url
		modify_mu.Unlock()
		resp, err := client.Get(url + "/healthy")
		modify_mu.Lock()
		unhealthy := err != nil || (resp != nil && resp.StatusCode != 200)
		if unhealthy {
			node := nodes_map[i]
			node.Healthy = false
			nodes_map[i] = node
			updateActiveNodes(i)
		}
		modify_mu.Unlock()
		if err == nil {
			resp.Body.Close()
		}
	}
}

func updateActiveNodes(key string) {
	for i := range active_nodes {
		if active_nodes[i].Id == key {
			log.Printf("node %d removed from active nodes", i)
			active_nodes[i] = active_nodes[len(active_nodes)-1]
			active_nodes = active_nodes[:len(active_nodes)-1]
			updateAlgorithms()
			break
		}
	}
}
