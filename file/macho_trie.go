package file

// The exports trie generator. A trie node is a ULEB terminal size (with
// the terminal data when nonzero: a ULEB flags word - always 0 here, plain
// absolute exports - and a ULEB image offset from the __TEXT vmaddr), then
// a ULEB child count and that many edges (a null-terminated string plus a
// ULEB child offset into the trie).
//
// The arrangement of the nodes in the byte stream reproduces the reference
// ld file: the root first, a fixed 4-byte dead zone, then the terminal
// leaves in name order, then the branch nodes deepest-first; the total is
// padded to 8. Child edges keep insertion order - the names arrive as
// __mh_execute_header plus the sorted globals - which is what makes the
// reference bytes (and any other prefix set) deterministic.

import (
	"slices"
	"strings"
)

// machoTrieNode - one node of the prefix trie, with serialization
// bookkeeping (name, depth, assigned offset and size).
type machoTrieNode struct {
	termSet bool
	term    uint64
	edges   []*machoTrieEdge

	name  string
	depth int
	off   int
	size  int
}

// machoTrieEdge - one child edge: an arbitrary-length string (chains are
// compressed into edges, so only real branch points become nodes).
type machoTrieEdge struct {
	s string
	n *machoTrieNode
}

// machoTrieInsert adds name with its terminal value, splitting edges where
// prefixes are shared.
func machoTrieInsert(root *machoTrieNode, name string, value uint64) {
	node := root
	rest := name

	for rest != "" {
		var match *machoTrieEdge
		best := 0

		for _, e := range node.edges {
			if c := machoTrieCommon(e.s, rest); c > best {
				best = c
				match = e
			}
		}

		if match == nil {
			leaf := &machoTrieNode{termSet: true, term: value, name: name}
			node.edges = append(node.edges, &machoTrieEdge{s: rest, n: leaf})
			return
		}

		if best == len(match.s) {
			node = match.n
			rest = rest[best:]
			continue
		}

		// split the edge: a mid node takes the shared prefix
		mid := &machoTrieNode{}
		mid.edges = append(mid.edges, &machoTrieEdge{s: match.s[best:], n: match.n})
		match.s = match.s[:best]
		match.n = mid

		if best == len(rest) {
			mid.termSet, mid.term, mid.name = true, value, name
			return
		}

		leaf := &machoTrieNode{termSet: true, term: value, name: name}
		mid.edges = append(mid.edges, &machoTrieEdge{s: rest[best:], n: leaf})
		return
	}

	node.termSet, node.term, node.name = true, value, name
}

// machoTrieCommon - the length of the longest common prefix of a and b.
func machoTrieCommon(a, b string) int {
	n := min(len(a), len(b))
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return i
		}
	}

	return n
}

// machoNewTrie builds the trie for names (already in the wanted edge
// order) and serializes it.
func machoNewTrie(names []string, valueOf func(string) uint64) []byte {
	root := &machoTrieNode{}
	for _, name := range names {
		machoTrieInsert(root, name, valueOf(name))
	}

	var leaves, branches []*machoTrieNode
	machoTrieWalk(root, 0, "", &leaves, &branches)
	slices.SortFunc(leaves, func(a, b *machoTrieNode) int {
		return strings.Compare(a.name, b.name)
	})
	slices.SortFunc(branches, func(a, b *machoTrieNode) int {
		if a.depth != b.depth {
			return b.depth - a.depth // deepest first
		}
		return strings.Compare(a.name, b.name)
	})

	// sizes depend on the ULEB lengths of the child offsets, which depend
	// on the offsets: iterate from the minimal seed to a fixpoint
	for changed := true; changed; {
		changed = false

		for _, n := range machoTrieNodes(root, leaves, branches) {
			if n.size != machoTrieSize(n) {
				n.size = machoTrieSize(n)
				changed = true
			}
		}

		cur := root.size + 4 // the dead zone of the reference layout
		for _, n := range leaves {
			n.off = cur
			cur += n.size
		}
		for _, n := range branches {
			n.off = cur
			cur += n.size
		}
	}

	out := machoTrieBytes(root)
	out = append(out, 0, 0, 0, 0)
	for _, n := range leaves {
		out = append(out, machoTrieBytes(n)...)
	}
	for _, n := range branches {
		out = append(out, machoTrieBytes(n)...)
	}

	return out[:pad8(len(out))]
}

// machoTrieWalk collects the non-root nodes with their names and depths.
func machoTrieWalk(
	n *machoTrieNode,
	depth int,
	prefix string,
	leaves, branches *[]*machoTrieNode,
) {
	n.name = prefix
	n.depth = depth

	for _, e := range n.edges {
		machoTrieWalk(e.n, depth+1, prefix+e.s, leaves, branches)
	}

	if depth == 0 {
		return // the root is serialized separately, first
	}

	if len(n.edges) == 0 {
		*leaves = append(*leaves, n)
	} else {
		*branches = append(*branches, n)
	}
}

// machoTrieNodes - root, leaves, branches (the serialization order).
func machoTrieNodes(root *machoTrieNode, leaves, branches []*machoTrieNode) []*machoTrieNode {
	out := make([]*machoTrieNode, 0, 1+len(leaves)+len(branches))
	out = append(out, root)
	out = append(out, leaves...)
	out = append(out, branches...)
	return out
}

// machoTrieSize - the serialized size of a node at its current offsets.
func machoTrieSize(n *machoTrieNode) int {
	size := 0

	if n.termSet {
		size += len(uleb(uint64(1 + len(uleb(n.term))))) // terminal size
		size += 1                                        // flags: a plain export
		size += len(uleb(n.term))
	} else {
		size += 1
	}

	size += len(uleb(uint64(len(n.edges))))
	for _, e := range n.edges {
		size += len(e.s) + 1 + len(uleb(uint64(e.n.off)))
	}

	return size
}

// machoTrieBytes - the node serialization at the frozen offsets.
func machoTrieBytes(n *machoTrieNode) []byte {
	out := []byte{}

	if n.termSet {
		out = append(out, uleb(uint64(1+len(uleb(n.term))))...)
		out = append(out, 0) // flags
		out = append(out, uleb(n.term)...)
	} else {
		out = append(out, 0)
	}

	out = append(out, uleb(uint64(len(n.edges)))...)
	for _, e := range n.edges {
		out = append(out, e.s...)
		out = append(out, 0)
		out = append(out, uleb(uint64(e.n.off))...)
	}

	return out
}
