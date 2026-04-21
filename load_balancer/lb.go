package lb

import (
	"encoding/json"
	"loadbalancer/algorithms"
	"loadbalancer/models"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync"
)

var nodes_map = make(map[string]models.Node)
var active_nodes []models.Node
var modify_mu sync.Mutex

var rr = algorithms.RoundRobin(len(active_nodes))
var wrr = algorithms.WRoundRobin(len(active_nodes), active_nodes)

func StartServer() {

	http.HandleFunc("/", lbHandler)
	http.HandleFunc("/register", registerNodeHandler)
	go StartHealthChecking()
	http.ListenAndServe(":8080", nil)

}
func updateAlgorithms() {
	rr = algorithms.RoundRobin(len(active_nodes))
	wrr = algorithms.WRoundRobin(len(active_nodes), active_nodes)
}
func registerNodeHandler(w http.ResponseWriter, r *http.Request) {
	modify_mu.Lock()
	defer modify_mu.Unlock()
	var node_info models.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&node_info)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	n, exist := nodes_map[node_info.Id]
	if exist {
		n.Healthy = true
		nodes_map[node_info.Id] = n
		log.Printf("node %s is back", node_info.Id)
		return
	}
	var node models.Node
	node.Initialize(node_info.Url, node_info.Weight, node_info.Id)
	nodes_map[node_info.Id] = node
	active_nodes = append(active_nodes, node)
	log.Printf("node %d is registered", len(active_nodes)-1)
	updateAlgorithms()
}
func lbHandler(w http.ResponseWriter, r *http.Request) {

	modify_mu.Lock()
	if len(active_nodes) == 0 {
		modify_mu.Unlock()
		http.Error(w, " no available nodes", http.StatusServiceUnavailable)
		return
	}
	modify_mu.Unlock()

	client_ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, "invalid client ip", http.StatusBadRequest)
		log.Printf("rejected request invalid client ip")
		return
	}
	if !checkLimit(client_ip) {
		http.Error(w, "too many requests", http.StatusTooManyRequests)
		log.Printf("rejected request , too many requests from client :%s", client_ip)
		return
	}
	algo_id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "invalid algorithm id", http.StatusBadRequest)
		log.Printf("error happend when getting algo_id %v", err)
		return
	}
	modify_mu.Lock()
	var node_id int
	used_lc := false
	switch algo_id {
	case 1:
		node_id = rr()
	case 2:
		node_id = wrr()
	case 3:
		node_id = algorithms.RandomLb(len(active_nodes))
	case 4:
		node_id = algorithms.LeastConnections(active_nodes)
		node := active_nodes[node_id]
		node.Connections++
		active_nodes[node_id] = node
		used_lc = true
	case 5:
		node_id = int(algorithms.HashLb(len(active_nodes), client_ip))
	}

	node_url, err := url.Parse(active_nodes[node_id].Url)
	if err != nil {
		modify_mu.Unlock()
		http.Error(w, "invalid node url", http.StatusInternalServerError)
		log.Printf("invalid node url parsing for node:%d", node_id)
		return
	}
	modify_mu.Unlock()

	proxy := httputil.NewSingleHostReverseProxy(node_url)
	proxy.ServeHTTP(w, r)
	log.Printf("request from : %s to node : %d (with algorithm : %d)", client_ip, node_id, algo_id)

	modify_mu.Lock()
	if used_lc {
		node := active_nodes[node_id]
		node.Connections--
		active_nodes[node_id] = node
	}
	modify_mu.Unlock()

}
