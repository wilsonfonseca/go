package paths

import (
	"fmt"
)

// Find finds all acyclic paths from source asset to destinations with length up to
// `MaxTrades`. Then it checks if the paths found satisfy destAmount.
func (e *Exchange) Find(source Asset, destination Asset) {
	sourceNode := e.GetNode(source)

	destinations := make(map[*Node]bool)
	destinations[e.GetNode(destination)] = true

	// Find all paths
	paths := e.dfs(sourceNode, destinations, 0, NewPath())

	// Check if found paths satisfy destAmount
	for _, path := range paths {
		path.Print()
	}
	fmt.Println(len(paths))
}

// dfs is the depth-first search implementation in the graph of assets and markets between them.
// There are however 2 modifications from the original DFS algorithm:
//   * Searching will stop when the path length is higher than `MaxTrades`
//   * Path with cycles are not allowed.
//
// Time complexity - simple analysis:
//
// Because we stop the search when depths is equal `MaxTrades`+1 paths will be split into
// `MaxTrades`+1 levels:
//   * From level 0 we will reach nodes that `source` asset has offers with.
//   * From level 1 we will reach nodes that adjacent nodes of `source` has offers with.
//   * ...
// Suppose that, on average, each asset is in 5 markets. The upper boundaries of dfs calls
// are the following:
//   * Level 0: 5 dfs calls.   = 5
//   * Level 1: 5^2 dfs calls. = 25
//   * Level 2: 5^3 dfs calls. = 124
//   * Level 3: 5^4 dfs calls. = 625
//   * Level 4: 5^5 dfs calls. = 3125
//                         sum = 3904
// However there are popular assets that other assets have offers with (ex. XLM, EURT, PHP)
// and because we do not allow cycles it means that every time a popular asset has been visited
// and a new less popular asset has offers with it it won't be visited. What is more, assets that
// are less popular are usually connected with only 1 popular asset. This makes the number of dfs
// calls even smaller.
//
// In the worst case scenario we can switch to BFS to find the shortest paths first and stop after
// checking X nodes.
func (e *Exchange) dfs(source *Node, destinations map[*Node]bool, depth int, p *path) []*path {
	pathsFound := []*path{}

	if depth > MaxTrades {
		// Path would be too long
		return pathsFound
	}

	p.Append(source)

	// Source is one of the nodes we're looking for
	if destinations[source] {
		return []*path{p}
	}

	for nextNode, _ := range source.Markets {
		if !p.IsVisited(nextNode) {
			newPaths := e.dfs(nextNode, destinations, depth+1, p.Clone())
			pathsFound = append(pathsFound, newPaths...)
		}
	}

	return pathsFound
}

func (e *Exchange) GetNode(asset Asset) *Node {
	node := e.Nodes[asset]
	if node != nil {
		return node
	}

	e.Nodes[asset] = &Node{
		Selling: asset,
		Markets: make(map[*Node][]Offer),
	}

	return e.Nodes[asset]
}
