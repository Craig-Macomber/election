package merkle

import (
    "bytes"
    "fmt"
)

// Return a [][]byte needed to prove the inclusion of the item at the passed index
// The payload of the item at index is the first value in the proof
func (t *Tree) InclusionProof(index int) [][]byte {
	if uint64(index) >= t.count {
		panic("Invalid index: too large")
	}
	if index < 0 {
		panic("Invalid index: negative")
	}
	h := GetHeight(t.count)
	return proveNode(h, t.root, index)
}

func proveNode(height int, n *node, index int) [][]byte {
	if height == 1 {
		if index != 0 {
			panic("Invalid index: non 0 for final node")
		}
		return [][]byte{n.label}
	}
	childIndex := index >> uint(height-2)
	nextIndex := index & (^(1<<uint(height-2)))
	b := proveNode(height-1, n.children[childIndex], nextIndex)
	otherChildIndex := (childIndex + 1) % 2
	if n.children[otherChildIndex]!=nil {
	    b = append(b, n.children[otherChildIndex].label)
	}
	return b
}

// The Root of a merkle tree for a client that does not store the tree
type Root struct {
    Count uint64
    Base []byte
}

// Proves the inclusion of an element at the given index with the value thats the first entry in proof
func (r *Root) CheckProof(proof [][]byte, index int) bool{
    h := GetHeight(r.Count)
    root,ok:=checkNode(h,proof,uint64(index),r.Count)
    base:=rootHash(r.Count,root)
    return ok && bytes.Equal(r.Base,base)
}

func checkNode(height int, proof [][]byte, index, count uint64) ([]byte,bool) {
	if len(proof)==0 {
	    fmt.Println("Empty")
	    return nil,false
	}
	if count<=index {
	    fmt.Println("bad count",count,index)
	    return nil,false
	}
	
	if height == 1 {
		if index != 0 || len(proof)!=1 {
		    fmt.Println("BAD",index,proof)
			return nil,false
		}
		return proof[0],true
	}

	childIndex := index >> uint(height-2)
	mask:=uint64(^(1<<uint(height-2)))
	nextIndex := index & mask
	
	var data []byte
	var ok bool
	
	h:=hashType.New()
	var nextCount uint64
	last:=len(proof)-1
	if childIndex==1 {
	    nextCount=count & mask
	    h.Write(proof[last])
	    data,ok=checkNode(height-1,proof[:last],nextIndex,nextCount)
	    h.Write(data)
	} else {
	    nextCount=count
	    if count>^mask {
	        nextCount=^mask
	    }
	    if count == nextCount {
	        data,ok=checkNode(height-1,proof,nextIndex,nextCount)
	        h.Write(data)
	    } else {
	        data,ok=checkNode(height-1,proof[:last],nextIndex,nextCount)
	        h.Write(data)
	        h.Write(proof[last])
	    }
	}
	
	hash:=h.Sum(make([]byte, 0))
	return hash,ok
}
