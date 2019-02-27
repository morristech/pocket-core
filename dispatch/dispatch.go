package dispatch

import (
	"math/big"
	"sort"

	"github.com/pokt-network/pocket-core/node"
)

type Dispatch struct {
	DevID       string            `json:"devid"`
	Blockchains []node.Blockchain `json:"blockchains"`
}

// NOTE: this call has been augmented for the Pocket Core MVP Centralized Dispatcher
// TODO see if this can be done more efficiently
// "Serve" formats Dispatch PL for an API request.
func Serve(dispatch *Dispatch) []byte {

}

/*
"Find" orders the nodes from smallest proximity from sessionKey to largest proximity to sessionKey.
// TODO convert to P2P -> currently just searches the peerlist
*/
func Find(sessionKey string) []node.Node {
	// create new key
	bigSessionKey := new(big.Int)
	bigSessionKey.SetString(sessionKey, 16)
	peerList := node.PeerList()
	peerList.Mux.Lock()
	defer peerList.Mux.Unlock()
	// map the nodes to the corresponding difference
	m := make(map[uint64]node.Node)
	// store the keys (to easily sort)
	keys := make([]uint64, len(peerList.M))
	// resulting array that holds the sorted nodes ordered by difference
	sortedNodes := make([]node.Node, len(peerList.M))
	var i = 0
	for gid, n := range peerList.M {
		// setup a new big integer to hold the converted ID
		id := new(big.Int)
		// convert the hex GID into a bigInteger for comparison
		id.SetString(gid.(string), 16)
		difference := big.NewInt(0).Sub(bigSessionKey, id)
		// take absolute of the difference for comparison
		difference.Abs(difference)
		m[difference.Uint64()] = n.(node.Node)
		keys[i] = difference.Uint64()
		i++
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for i, k := range keys {
		sortedNodes[i] = m[k]
	}
	return sortedNodes
}
