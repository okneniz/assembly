// Package dtree is a decision tree over mask rules: given a word, it selects
// the first matching rule (first-match-wins) without scanning the list
// linearly.
//
// It is built from rules [{Mask, Match, Payload}]: a node splits the group by
// bits defined by the masks of all its rules that distinguish them; if there
// are no such bits, by a single bit covered by some of the rules (rules
// without the bit in their mask go into both branches, the original order is
// preserved); rules not distinguishable by any bit remain an ordered list in
// a leaf. Thus the semantics of the list (order = priority, intersections
// allowed) is preserved unchanged.
//
// Experiment: a candidate for moving to parsec - "mask-based parsing" is
// broader than a single ISA (protocol headers, opcode dispatchers).
package dtree

// Rule is a mask rule: a word belongs to it if (word & Mask) == Match.
// Payload is what Lookup returns (constructor, mnemonic, ...); for masking
// entries (matched, but there is no construction) a zero Payload is stored.
type Rule[T any] struct {
	Mask    uint32
	Match   uint32
	Payload T
}

// Tree is a decision tree; built by New, safe for reading from goroutines.
type Tree[T any] struct {
	root *node[T]
}

// node is an internal node (branches != nil, mask is the splitting bits) or
// a leaf (rules in the original order).
type node[T any] struct {
	mask     uint32
	branches map[uint32]*node[T]
	rules    []Rule[T]
}

// New builds the tree; the order of rules defines priority (the first match
// wins), the order within is preserved by construction.
func New[T any](rules []Rule[T]) *Tree[T] {
	return &Tree[T]{root: build(rules)}
}

// Lookup returns the Payload of the first rule (in New order) that word
// belongs to; ok=false means no match.
func (t *Tree[T]) Lookup(word uint32) (payload T, ok bool) {
	n := t.root
	for n != nil && n.branches != nil {
		n = n.branches[word&n.mask]
	}

	if n == nil {
		return payload, false
	}

	for _, r := range n.rules {
		if word&r.Mask == r.Match {
			return r.Payload, true
		}
	}

	return payload, false
}

// build splits the group recursively for as long as possible; a leaf is what
// remains.
func build[T any](rules []Rule[T]) *node[T] {
	if len(rules) <= 1 {
		return &node[T]{rules: rules}
	}

	// Bits defined by the masks of all rules of the group...
	covered := ^uint32(0)
	for _, r := range rules {
		covered &= r.Mask
	}

	// ... of them, those distinguishing the values (split without duplication).
	var orAll uint32
	andAll := ^uint32(0)
	for _, r := range rules {
		orAll |= r.Match & covered
		andAll &= r.Match & covered
	}

	if vary := orAll &^ andAll; vary != 0 {
		groups := make(map[uint32][]Rule[T])
		for _, r := range rules {
			k := r.Match & vary
			groups[k] = append(groups[k], r)
		}

		return fork(vary, groups)
	}

	// Partially covered bit: rules without the bit in their mask go into both
	// branches (duplication). The bit must distinguish the covered part (both
	// sides are strictly smaller than the group). Duplication is costly: a set
	// with broad "indifferent" rules (catch-alls) blows the tree up
	// exponentially, so we split by the bit with the fewest duplicates and
	// only if they are no more than a quarter of the group; otherwise, a leaf
	// list.
	best, bestDup := -1, 0
	for b := range uint(32) {
		bit := uint32(1) << b

		var hasZero, hasOne bool
		dup := 0
		for _, r := range rules {
			switch {
			case r.Mask&bit == 0:
				dup++
			case r.Match&bit == 0:
				hasZero = true
			default:
				hasOne = true
			}
		}

		if !hasZero || !hasOne {
			continue // the bit does not distinguish the covered part of the group
		}

		if best < 0 || dup < bestDup {
			best, bestDup = int(b), dup
		}
	}

	if best < 0 || bestDup > len(rules)/4 {
		// no splitting bit without costly duplication (intersections,
		// indifference) - a priority-ordered leaf list
		return &node[T]{rules: rules}
	}

	bit := uint32(1) << uint(best)

	var zero, one []Rule[T]
	for _, r := range rules {
		switch {
		case r.Mask&bit == 0: // indifferent - into both branches
			zero = append(zero, r)
			one = append(one, r)
		case r.Match&bit == 0:
			zero = append(zero, r)
		default:
			one = append(one, r)
		}
	}

	// Branch keys are field values: match & bit (0 or bit), the way Lookup
	// searches (word & mask).
	return fork(bit, map[uint32][]Rule[T]{0: zero, bit: one})
}

// fork assembles a node with branches and descends into each group
// recursively.
func fork[T any](mask uint32, groups map[uint32][]Rule[T]) *node[T] {
	n := &node[T]{
		mask:     mask,
		branches: make(map[uint32]*node[T], len(groups)),
	}

	for k, g := range groups {
		n.branches[k] = build(g)
	}

	return n
}

// each traverses the tree depth-first, passing the node and its depth (the
// root is 0).
func (t *Tree[T]) each(visit func(n *node[T], depth int)) {
	if t.root == nil {
		return
	}

	var walk func(n *node[T], depth int)
	walk = func(n *node[T], depth int) {
		visit(n, depth)

		for _, c := range n.branches {
			walk(c, depth+1)
		}
	}

	walk(t.root, 0)
}

// Nodes is the number of tree nodes (internal + leaves); for memory
// estimation.
func (t *Tree[T]) Nodes() int {
	count := 0
	t.each(func(_ *node[T], _ int) { count++ })

	return count
}

// Leaves is the number of leaves (terminal rule lists).
func (t *Tree[T]) Leaves() int {
	count := 0
	t.each(func(n *node[T], _ int) {
		if n.branches == nil {
			count++
		}
	})

	return count
}

// MaxDepth is the Lookup descent depth (<= 32: a node consumes at least one
// bit).
func (t *Tree[T]) MaxDepth() int {
	depth := 0
	t.each(func(_ *node[T], d int) {
		if d+1 > depth {
			depth = d + 1
		}
	})

	return depth
}
