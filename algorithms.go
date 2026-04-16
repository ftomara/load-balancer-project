package main

import (
	"hash/fnv"
	"math/rand"
)

func round_robin(nodes_count int) func() int {
	next := 0
	return func() int {
		n := next % nodes_count
		next++
		return n
	}
}

func w_round_robin(nodes_count int, nodes_weight []int) func() int {
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

func random_lb(nodes_count int) int {
	return rand.Intn(nodes_count)
}

func least_connections(connections []int) int {
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

func hash_lb(node_count int, ip string) uint32 {
	hasher := fnv.New32()
	hasher.Write([]byte(ip))
	node := hasher.Sum32() % uint32(node_count)
	return node
}
