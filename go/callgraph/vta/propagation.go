// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vta

import (
	"go/types"

	"github.com/kent0106/gotools/go/ssa"

	"github.com/kent0106/gotools/go/types/typeutil"
)

// scc computes strongly connected components (SCCs) of `g` using the
// classical Tarjan's algorithm for SCCs. The result is a pair <m, id>
// where m is a map from nodes to unique id of their SCC in the range
// [0, id). The SCCs are sorted in reverse topological order: for SCCs
// with ids X and Y s.t. X < Y, Y comes before X in the topological order.
func scc(g vtaGraph) (map[node]int, int) {
	// standard data structures used by Tarjan's algorithm.
	var index uint64
	var stack []node
	indexMap := make(map[node]uint64)
	lowLink := make(map[node]uint64)
	onStack := make(map[node]bool)

	nodeToSccID := make(map[node]int)
	sccID := 0

	var doSCC func(node)
	doSCC = func(n node) {
		indexMap[n] = index
		lowLink[n] = index
		index = index + 1
		onStack[n] = true
		stack = append(stack, n)

		for s := range g[n] {
			if _, ok := indexMap[s]; !ok {
				// Analyze successor s that has not been visited yet.
				doSCC(s)
				lowLink[n] = min(lowLink[n], lowLink[s])
			} else if onStack[s] {
				// The successor is on the stack, meaning it has to be
				// in the current SCC.
				lowLink[n] = min(lowLink[n], indexMap[s])
			}
		}

		// if n is a root node, pop the stack and generate a new SCC.
		if lowLink[n] == indexMap[n] {
			for {
				w := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				onStack[w] = false
				nodeToSccID[w] = sccID
				if w == n {
					break
				}
			}
			sccID++
		}
	}

	index = 0
	for n := range g {
		if _, ok := indexMap[n]; !ok {
			doSCC(n)
		}
	}

	return nodeToSccID, sccID
}

func min(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}

// propType represents type information being propagated
// over the vta graph. f != nil only for function nodes
// and nodes reachable from function nodes. There, we also
// remember the actual *ssa.Function in order to more
// precisely model higher-order flow.
type propType struct {
	typ types.Type
	f   *ssa.Function
}

// propTypeMap is an auxiliary structure that serves
// the role of a map from nodes to a set of propTypes.
type propTypeMap struct {
	nodeToScc  map[node]int
	sccToTypes map[int]map[propType]bool
}

// propTypes returns a set of propTypes associated with
// node `n`. If `n` is not in the map `ptm`, nil is returned.
//
// Note: for performance reasons, the returned set is a
// reference to existing set in the map `ptm`, so any updates
// to it will affect `ptm` as well.
func (ptm propTypeMap) propTypes(n node) map[propType]bool {
	id, ok := ptm.nodeToScc[n]
	if !ok {
		return nil
	}
	return ptm.sccToTypes[id]
}

// propagate reduces the `graph` based on its SCCs and
// then propagates type information through the reduced
// graph. The result is a map from nodes to a set of types
// and functions, stemming from higher-order data flow,
// reaching the node. `canon` is used for type uniqueness.
func propagate(graph vtaGraph, canon *typeutil.Map) propTypeMap {
	nodeToScc, sccID := scc(graph)

	// We also need the reverse map, from ids to SCCs.
	sccs := make(map[int][]node, sccID)
	for n, id := range nodeToScc {
		sccs[id] = append(sccs[id], n)
	}

	// Initialize sccToTypes to avoid repeated check
	// for initialization later.
	sccToTypes := make(map[int]map[propType]bool, sccID)
	for i := 0; i <= sccID; i++ {
		sccToTypes[i] = nodeTypes(sccs[i], canon)
	}

	for i := len(sccs) - 1; i >= 0; i-- {
		nextSccs := make(map[int]struct{})
		for _, node := range sccs[i] {
			for succ := range graph[node] {
				nextSccs[nodeToScc[succ]] = struct{}{}
			}
		}
		// Propagate types to all successor SCCs.
		for nextScc := range nextSccs {
			mergeTypes(sccToTypes[nextScc], sccToTypes[i])
		}
	}
	return propTypeMap{nodeToScc: nodeToScc, sccToTypes: sccToTypes}
}

// nodeTypes returns a set of propTypes for `nodes`. These are the
// propTypes stemming from the type of each node in `nodes` plus.
func nodeTypes(nodes []node, canon *typeutil.Map) map[propType]bool {
	types := make(map[propType]bool)
	for _, n := range nodes {
		if hasInitialTypes(n) {
			types[getPropType(n, canon)] = true
		}
	}
	return types
}

// getPropType creates a propType for `node` based on its type.
// propType.typ is always node.Type(). If node is function, then
// propType.val is the underlying function; nil otherwise.
func getPropType(node node, canon *typeutil.Map) propType {
	t := canonicalize(node.Type(), canon)
	if fn, ok := node.(function); ok {
		return propType{f: fn.f, typ: t}
	}
	return propType{f: nil, typ: t}
}

// mergeTypes merges propTypes in `rhs` to `lhs`.
func mergeTypes(lhs, rhs map[propType]bool) {
	for typ := range rhs {
		lhs[typ] = true
	}
}
