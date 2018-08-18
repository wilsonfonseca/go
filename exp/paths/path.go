package paths

import (
	"fmt"
)

func NewPath() *path {
	return &path{
		Nodes:   make([]*Node, 0, MaxTrades+1),
		Visited: make(map[*Node]bool),
	}
}

func (p *path) Append(n *Node) {
	p.Nodes = append(p.Nodes, n)
	p.Visited[n] = true
}

func (p *path) IsVisited(n *Node) bool {
	return p.Visited[n]
}

func (p *path) Clone() *path {
	newPath := NewPath()

	for _, n := range p.Nodes {
		newPath.Append(n)
	}

	return newPath
}

func (p *path) Print() {
	fmt.Println("Length", len(p.Nodes))
	for _, node := range p.Nodes {
		fmt.Printf("%+v -> ", node.Selling.String())
	}
	fmt.Print("\n\n")
}
