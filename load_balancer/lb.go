package lb

import (
	"fmt"
	"loadbalancer/algorithms"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync"
)

var nodes = []string{"http://node1:8081", "http://node2:8082", "http://node3:8083", "http://node4:8084", "http://node5:8085"}
var nodes_weight = []int{3, 5, 1, 3, 8}
var nodes_connections = []int{0, 0, 0, 0, 0}
var lc_mu sync.Mutex

var rr = algorithms.RoundRobin(len(nodes))
var wrr = algorithms.WRoundRobin(len(nodes), nodes_weight)

func StartServer() {

	http.HandleFunc("/", lbHandler)
	http.ListenAndServe(":8080", nil)

}

func lbHandler(w http.ResponseWriter, r *http.Request) {

	algo_id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		fmt.Printf("error happend when getting algo_id %v", err)
		http.Error(w, "invalid algorithm id", http.StatusBadRequest)
		return
	}

	client_ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, "invalid client ip", http.StatusBadRequest)
		return
	}

	var node_id int
	used_lc := false
	switch algo_id {
	case 1:
		node_id = rr()
	case 2:
		node_id = wrr()
	case 3:
		node_id = algorithms.RandomLb(len(nodes))
	case 4:
		lc_mu.Lock()
		node_id = algorithms.LeastConnections(nodes_connections)
		nodes_connections[node_id]++
		used_lc = true
		lc_mu.Unlock()
	case 5:
		node_id = int(algorithms.HashLb(len(nodes), client_ip))
	}

	node_url, err := url.Parse(nodes[node_id])
	if err != nil {
		http.Error(w, "invalid node url", http.StatusInternalServerError)
		return
	}
	fmt.Printf("request is going to node : %v",node_id)
	proxy := httputil.NewSingleHostReverseProxy(node_url)
	proxy.ServeHTTP(w, r)
	lc_mu.Lock()
	if used_lc {
		nodes_connections[node_id]--
	}
	lc_mu.Unlock()

}
