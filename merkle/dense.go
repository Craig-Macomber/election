// Package merkle provides a binary merkle tree implementation.
package merkle

import (
	"bytes"
	"crypto"
	_ "crypto/sha256"
	"fmt"
	"encoding/binary"
)

var hashType crypto.Hash = crypto.SHA256

// A merkle tree for a user that stores the entire tree
// Specifically this tree is left a leaning balanced binary tree
// Where each node holds the hash of its leaves
// And the rootHash is the root node hashed with the count
// This tree is immutable
type Tree struct {
	count uint64
	root  *node
	rootHash []byte
}


type node struct {
	label  []byte
	children [2]*node // if all nil, leaf node

	// Representation invariants:
	// if children[0] is nil, children[1] is nil
	// if both children non nil:
	//      label is hash of (children[0].label + children[1].label)
	// if leaf: label is arbitrary data
	// else if children[1] is nil, label=hash(children[0].label)
}

func (t Tree) Count() uint64 {
	return t.count
}

// The hash/root of an empty tree does not matter
func (t Tree) Root() []byte {
	return t.rootHash
}

// All trees should pass Validate, unless they are invalid, which should only happen
// if incorrectly built or modified.
// Checks the rep invariants
func (t Tree) Validate() error {
	count, height, error := t.root.validate()
	if error != nil {
		return error
	}
	if count != t.count {
		return fmt.Errorf("Incorrect count. Was %d, should be %d", t.count, count)
	}
	if height != GetHeight(count) {
		return fmt.Errorf("Incorrect height. Was %d, should be %d", height, GetHeight(count))
	}
	
	rootLabel:=make([]byte,0)
	if height>0 {
	    rootLabel=t.root.label
	}
	h:=rootHash(count,rootLabel)
	if !bytes.Equal(t.rootHash,h) {
	    return fmt.Errorf("Incorrect rootHash")
	}
	return nil
}

// Checks the rep invariants
func (t *node) validate() (count uint64, height int, err error) {
	if t == nil {
		return 0, 0, nil
	}
	if t.children[0] == nil {
		if t.children[1] != nil {
			return 0, 0, fmt.Errorf("Invalid Node: Node missing first child, but has second")
		}
		// Leaf node
		return 1, 1, nil
	}

	// Not a leaf node
	count, height, err = t.children[0].validate()
	if err != nil {
		return
	}
	if t.children[1] != nil {
		count2, height2, err2 := t.children[1].validate()
		count += count2
		if err2 != nil {
			return count, height, err2
		}
		if height2 != height {
			return count, height, fmt.Errorf("Invalid Node: height mismatch between children")
		}
	}
	h := makeHash(t.children[0], t.children[1])
	if !bytes.Equal(h, t.label) {
		return 0, 0, fmt.Errorf("Invalid Node: Node hash mismatch")
	}

	height += 1
	return
}

func rootHash(count uint64,data []byte) []byte {
	h := hashType.New()
    h.Write(data)
    binary.Write(h,binary.LittleEndian,count)
	return h.Sum(make([]byte, 0))
}

func makeHash(left, right *node) []byte {
	h := hashType.New()
	if left != nil {
		h.Write(left.label)
		if right != nil {
			h.Write(right.label)
		}
	}
	return h.Sum(make([]byte, 0))
}

// Returns the height of the tree containing count leaf nodes.
// This the number of nodes (including the final leaf) from the root to
// any leaf.
func GetHeight(count uint64) int {
	if count == 0 {
		return 0
	}
	height := 0
	for count > (1 << uint(height)) {
		height++
	}
	return height + 1
}

// Build a tree
func Build(data [][]byte) *Tree {
	count := uint64(len(data))
	height := GetHeight(count)
	node, leftOverData := buildNode(data, height)
	if len(leftOverData) != 0 {
		panic("Build failed to consume all data")
	}
	rootLabel:=make([]byte,0)
	if height>0 {
	    rootLabel=node.label
	}
	hash:=rootHash(count,rootLabel)
	t := Tree{count, node, hash}
	return &t
}

// returns a node and the left over data not used by it
func buildNode(data [][]byte, height int) (*node, [][]byte) {
	if height == 0 || len(data) == 0 {
		return nil, data
	}
	if height == 1 {
		// leaf
		return &node{label: data[0]}, data[1:]
	}
	n0, data := buildNode(data, height-1)
	n1, data := buildNode(data, height-1)
	
	hash := makeHash(n0, n1)
	return &node{label: hash, children: [2]*node{n0, n1}}, data
}
