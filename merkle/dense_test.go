package merkle

import (
	"testing"
	"fmt"
)

func TestGetHeight(t *testing.T) {
	data := [][2]int{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 3},
		{4, 3},
		{255, 9},
		{256, 9},
		{257, 10},
	}
	for _, v := range data {
		h := GetHeight(uint64(v[0]))
		if !(v[1] == h) {
			t.Errorf("GetHeight(%d)!=%d (was %d)", v[0], v[1], h)
		}
	}
}

func TestGetHeight2(t *testing.T) {
	for i := 1; i < 1000; i++ {
		h := GetHeight(uint64(i))
		upperBound := 1 << uint(h-1)
		lowerBound := (1 << uint(h-2)) + 1
		if i < lowerBound {
			t.Errorf("GetHeight(%d) too high: %d", i, h)
		}
		if i > upperBound {
			t.Errorf("GetHeight(%d) too low: %d", i, h)
		}
	}
}

func TestBuild(t *testing.T) {
	data := [][]byte{}
	const count = 100
	for i := 0; i < count; i++ {
		d := [1]byte{byte(i)}
		data = append(data, d[:])
	}
	tree := Build(data)
	err := tree.Validate()
	if err != nil {
		t.Errorf("%s", err)
	}
	if tree.Count() != count {
		t.Errorf("GetCount() != %d (was %d)", count, tree.Count())
	}
	
	r:=Root{count,tree.Root()}
	
	for i := 0; i < count; i++ {
		p:=tree.InclusionProof(i)
		fmt.Println(p)
		ok:=r.CheckProof(p,i)
		if !ok {
		    t.Errorf("proof %d failed", i)
		}
	}
	t.Errorf("proof %d failed")
	// TODO: check wrong proofs fail
	
}
