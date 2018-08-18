package paths

import (
	"github.com/stellar/go/xdr"
)

// MaxTrades defines the maximum number of trades on paths to find.
// This should not be higher than `6` which is the max allowed by the
// path_payment operation.
const MaxTrades = 4

// WARNING: this is experimental and requires more work, testing and performance checks.
//
// Exchange represents the current state of the Stellar Distributed Exchange. It contains
// all active offers in the network and allows searching paths between assets.
//
// This object should be updated with the current state of the network every few seconds.
// Finding paths in memory does not generate a load on the core DB as only a single query
// that update the state is needed. Memory requirements for this object still need to be
// checked carefully.
//
// In the future the object state can be changed by checking the stream of operations/effects
// (add new offers / edit offers after trades / delete offers).
//
// TODO: add `AddOffer` function to construct a graph. It should check if offers are ordered
// by the price.
type Exchange struct {
	Nodes map[Asset]*Node
}

// Node represents a verticle (asset) in the graph of Stellar Distributed Exchange.
// Markets are edges in the graph between two assets and contain a list of offers
// selling `Node.Selling` asset and buying `Markets[*Node].Selling` asset.
// IMPORANT: Offers in each market must be ordered from the lowest price!
type Node struct {
	Selling Asset
	Amounts int64
	Markets map[*Node][]Offer
}

type Asset struct {
	Type   xdr.AssetType
	Code   string
	Issuer string
}

type Offer struct {
	Amount int64
	Pricen int32
	Priced int32
}

// path represents a path found using pathfinding algorithm.
type path struct {
	Nodes   []*Node
	Visited map[*Node]bool
}
