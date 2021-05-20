package my_merkleTree

import (
	//"encoding/hex"
	"errors"
	"github.com/miguelmota/go-solidity-sha3"
)

type Tree struct {
	head *subTree
	//hash []byte

	currentIndex uint64
	proofIndex   uint64
	proofSet     [][]byte
	proofTree    bool

	//this was added
	cachedTree bool
}

type subTree struct {
	next   *subTree
	height int
	sum    []byte
}

//
//leafSum returns the hash created from data inserted to form a leaf
//Leaf sums are calculate using :
//	Hash("0x00" || data)
func leafSum(data []byte) []byte {
	hash := solsha3.SoliditySHA3(
		[]string{"string", "string"},
		[]interface{}{
			"0x00",
			string(data),
		},
	)
	return []byte(hash)
}

//nodeSum returns the hash from two siblings nodes being combined
// with the type :
// Hash(0x01 || left sibling sum || right sibling sum)
func nodeSum(a, b []byte) []byte {
	hash := solsha3.SoliditySHA3(
		[]string{"string", "string", "string"},
		[]interface{}{
			"0x01",
			string(a),
			string(b),
		},
	)
	return []byte(hash)
}

//New creates a new Tree
func New() *Tree {
	return &Tree{}
}

//Prove creates a proof that the leaf at the established  index  is an
//element of the Merkle tree.
func (t *Tree) Prove() (merkleRoot []byte, proofSet [][]byte, proofIndex uint64, numLeaves uint64) {
	if !t.proofTree {
		panic("wrong usage : can't call prove on a tree if setIndex is not already called.")
	}

	//Return nill if the tree is empty or if the proof index has not yet been
	//reached.
	if t.head == nil || len(t.proofSet) == 0 {
		return t.Root(), nil, t.proofIndex, t.currentIndex
	}
	proofSet = t.proofSet

	current := t.head
	for current.next != nil && current.next.height < len(proofSet)-1 {
		current = joinSubTrees(current.next, current)
	}

	if current.next != nil && current.next.height == len(proofSet)-1 {
		proofSet = append(proofSet, current.sum)
		current = current.next
	}
	current = current.next
	for current != nil {
		proofSet = append(proofSet, current.sum)
		current = current.next
	}
	return t.Root(), proofSet, t.proofIndex, t.currentIndex
}

//Push will add data to the set, building out the Merkle Tree and Root.
//tree does not remember all elements being addd, instead only keeping
//the log(n) elements that are necessary to build a Merkle Root and
//keeping the log(n) elements necessary to build a proof that a piece of
//data is in the Merkle Tree.

func (t *Tree) Push(data []byte) {
	if t.currentIndex == t.proofIndex {
		t.proofSet = append(t.proofSet, data)
	}

	t.head = &subTree{
		next:   t.head,
		height: 0,
		//sum:    leafSum(data),
	}

	if t.cachedTree {
		t.head.sum = data
	} else {
		t.head.sum = leafSum(data)
	}

	//join subTrees if possible
	t.joinAllSubTrees()

	//update current index
	t.currentIndex++
}

func (t *Tree) PushSubTree(height int, sum []byte) error {
	newIndex := t.currentIndex + 1<<uint64(height)
	if t.proofTree && (t.currentIndex == t.proofIndex) ||
		(t.currentIndex < t.proofIndex && t.proofIndex < newIndex) {
		return errors.New("the cached tree  should not contain the element to prove")
	}

	if t.head != nil && height > t.head.height {
		return errors.New("can't add a subtree that is smallest than the smallest subtree")
	}

	t.head = &subTree{
		height: height,
		next:   t.head,
		sum:    sum,
	}

	t.joinAllSubTrees()

	t.currentIndex = newIndex
	return nil
}

//Root returns the Merkle root of the data that has been pushed.
func (t *Tree) Root() []byte {
	if t.head == nil {
		return nil
	}
	current := t.head
	for current.next != nil {
		current = joinSubTrees(current.next, current)

	}
	//Return a copy to prevent leaking a pointer to internal data
	return append(current.sum[:0:0], current.sum...)
}

//SetIndex will tell the Tree to create a storage proof for the leaf at
//the input index. SetIndex must be called on an empty tree
func (t *Tree) SetIndex(i uint64) error {
	if t.head != nil {
		return errors.New("cannot call SetIndex on Tree if Tree has not been reset")
	}
	t.proofTree = true
	t.proofIndex = i
	return nil
}

func (t *Tree) joinAllSubTrees() {
	for t.head.next != nil && t.head.height == t.head.next.height {

		if t.head.height == len(t.proofSet)-1 {
			leaves := uint64(1 << uint(t.head.height))
			mid := (t.currentIndex / leaves) * leaves

			if t.proofIndex < mid {
				t.proofSet = append(t.proofSet, t.head.sum)
			} else {
				t.proofSet = append(t.proofSet, t.head.next.sum)
			}
		}

		t.head = joinSubTrees(t.head.next, t.head)
	}
}

func joinSubTrees(a, b *subTree) *subTree {
	return &subTree{
		next:   a.next,
		height: a.height + 1,
		sum:    nodeSum(a.sum, b.sum),
	}

}
