package main

import (
	"os"
	"loadbalancer/node"
)

func main() {

	port := os.Getenv("PORT")
	node.StartServer(port)
}
