package algorithms

import (
	"hash/fnv"
	"math/rand"
)

func RoundRobin(nodes_count int) func() int {
	next := 0
	return func() int {
		n := next % nodes_count
		next++
		return n
	}
}

func WRoundRobin(nodes_count int, nodes_weight []int) func() int {
	next_node := 0
	current_node_limit := 0

	return func() int {
		n := next_node % nodes_count
		if current_node_limit < nodes_weight[n] {
			current_node_limit++
			return n
		} else {
			next_node++
			current_node_limit = 1
			n = next_node % nodes_count
			return n
		}
	}
}

func RandomLb(nodes_count int) int {
	return rand.Intn(nodes_count)
}

func LeastConnections(connections []int) int {
	mini := connections[0]
	min_index := 0

	for i := 0; i < len(connections); i++ {
		if connections[i] < mini {
			mini = connections[i]
			min_index = i
		}
	}
	return min_index
}

func HashLb(node_count int, ip string) uint32 {
	hasher := fnv.New32()
	hasher.Write([]byte(ip))
	node := hasher.Sum32() % uint32(node_count)
	return node
}
