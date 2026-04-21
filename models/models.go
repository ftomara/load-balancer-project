package models

type Node struct {
	Id          string
	Url         string
	Weight      int
	Connections int
	Healthy     bool
}
type RegisterRequest struct {
	Id     string    `json:"id"`
	Url    string `json:"url"`
	Weight int    `json:"weight"`
}

func (node *Node) Initialize(url string, weight int, id string) {
	node.Id = id
	node.Url = "http://" + url
	node.Weight = weight
	node.Connections = 0
	node.Healthy = true
}
