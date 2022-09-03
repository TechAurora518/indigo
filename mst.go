package mst

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"math/bits"

	"github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	cbor "github.com/ipfs/go-ipld-cbor"
	cbg "github.com/whyrusleeping/cbor-gen"
)

type MerkleSearchTree struct {
	fanout int
	cst    cbor.IpldStore

	entries  []NodeEntry
	layer    int
	pointer  cid.Cid
	validPtr bool

	cachedCid cid.Cid
}

func NewMST(cst cbor.IpldStore, fanout int, ptr cid.Cid, entries []NodeEntry, layer int) *MerkleSearchTree {
	mst := &MerkleSearchTree{
		cst:      cst,
		fanout:   fanout,
		pointer:  ptr,
		layer:    layer,
		entries:  entries,
		validPtr: ptr.Defined(),
	}

	return mst
}

func (mst *MerkleSearchTree) newTree(entries []NodeEntry) *MerkleSearchTree {
	return NewMST(mst.cst, mst.fanout, cid.Undef, entries, mst.layer)
}

func create(ctx context.Context, cst cbor.IpldStore, entries []NodeEntry, layer int, fanout int) (*MerkleSearchTree, error) {
	ptr, err := cidForEntries(ctx, entries, cst)
	if err != nil {
		return nil, err
	}

	return NewMST(cst, fanout, ptr, entries, layer), nil
}

func cborGet(ctx context.Context, bs blockstore.Blockstore, c cid.Cid, out cbg.CBORUnmarshaler) error {
	blk, err := bs.Get(ctx, c)
	if err != nil {
		return err
	}

	if err := out.UnmarshalCBOR(bytes.NewReader(blk.RawData())); err != nil {
		return err
	}

	return nil
}

const (
	EntryUndefined = 0
	EntryLeaf      = 1
	EntryTree      = 2
)

type NodeEntry struct {
	Kind int
	Key  []byte
	Val  cid.Cid
	Tree *MerkleSearchTree
}

func treeEntry(t *MerkleSearchTree) NodeEntry {
	return NodeEntry{
		Kind: EntryTree,
		Tree: t,
	}
}

func (ne NodeEntry) isTree() bool {
	return ne.Kind == EntryTree
}

func (ne NodeEntry) isLeaf() bool {
	return ne.Kind == EntryLeaf
}

func (ne NodeEntry) isUndefined() bool {
	return ne.Kind == EntryUndefined
}

type TreeEntry struct {
	P int64    `cborgen:"p"`
	K []byte   `cborgen:"k"`
	V cid.Cid  `cborgen:"v"`
	T *cid.Cid `cborgen:"t"`
}

type NodeData struct {
	L *cid.Cid    `cborgen:"l"`
	E []TreeEntry `cborgen:"e"`
}

func (mst *MerkleSearchTree) Add(ctx context.Context, key []byte, val cid.Cid, knownZeros int) (*MerkleSearchTree, error) {
	keyZeros := knownZeros // is this really just to avoid rerunning the leading zeros hash?
	if keyZeros < 0 {
		keyZeros = leadingZerosOnHash(key, mst.fanout)
	}

	layer, err := mst.getLayer(ctx)
	if err != nil {
		return nil, err
	}

	newLeaf := NodeEntry{
		Kind: EntryLeaf,
		Key:  key,
		Val:  val,
	}

	if keyZeros == layer {
		// it belongs to me
		index, err := mst.findGtOrEqualLeafIndex(ctx, key)
		if err != nil {
			return nil, err
		}

		found, err := mst.atIndex(index)
		if err != nil {
			return nil, err
		}

		if found.isLeaf() && bytes.Equal(found.Key, key) {
			return nil, fmt.Errorf("value already set at key: %x", key)
		}

		prevNode, err := mst.atIndex(index - 1)
		if err != nil {
			return nil, err
		}

		if prevNode.isUndefined() || prevNode.isLeaf() {
			return mst.spliceIn(ctx, newLeaf, index)
		}

		left, right, err := prevNode.Tree.splitAround(ctx, key)
		if err != nil {
			return nil, err
		}

		return mst.replaceWithSplit(ctx, index-1, left, newLeaf, right)
	} else if keyZeros > layer {
		index, err := mst.findGtOrEqualLeafIndex(ctx, key)
		if err != nil {
			return nil, err
		}

		prevNode, err := mst.atIndex(index - 1)
		if err != nil {
			return nil, err
		}

		if !prevNode.isUndefined() && prevNode.isTree() {
			newSubtree, err := prevNode.Tree.Add(ctx, key, val, keyZeros)
			if err != nil {
				return nil, err
			}

			return mst.updateEntry(ctx, index-1, NodeEntry{
				Kind: EntryTree,
				Tree: newSubtree,
			})
		} else {
			subTree, err := mst.createChild(ctx)
			if err != nil {
				return nil, err
			}

			newSubTree, err := subTree.Add(ctx, key, val, keyZeros)
			if err != nil {
				return nil, err
			}

			return mst.spliceIn(ctx, treeEntry(newSubTree), index)
		}
	} else {
		left, right, err := mst.splitAround(ctx, key)
		if err != nil {
			return nil, err
		}

		layer, err := mst.getLayer(ctx)
		if err != nil {
			return nil, err
		}

		extraLayersToAdd := keyZeros - layer

		for i := 1; i < extraLayersToAdd; i++ {
			// seems bad if both left and right are non nil
			if left != nil {
				par, err := left.createParent(ctx)
				if err != nil {
					return nil, err
				}
				left = par
			}

			if right != nil {
				par, err := right.createParent(ctx)
				if err != nil {
					return nil, err
				}
				right = par
			}

		}
		var updated []NodeEntry
		if left != nil {
			updated = append(updated, treeEntry(left))
		}
		if right != nil {
			updated = append(updated, treeEntry(right))
		}

		newRoot, err := create(ctx, mst.cst, updated, keyZeros, mst.fanout)
		if err != nil {
			return nil, err
		}
		// why invalidate?
		newRoot.validPtr = false

		return newRoot, nil
	}
}

func (mst *MerkleSearchTree) createParent(ctx context.Context) (*MerkleSearchTree, error) {
	layer, err := mst.getLayer(ctx)
	if err != nil {
		return nil, err
	}

	return NewMST(mst.cst, mst.fanout, cid.Undef, []NodeEntry{treeEntry(mst)}, layer+1), nil
}

func (mst *MerkleSearchTree) createChild(ctx context.Context) (*MerkleSearchTree, error) {
	layer, err := mst.getLayer(ctx)
	if err != nil {
		return nil, err
	}

	return NewMST(mst.cst, mst.fanout, cid.Undef, nil, layer-1), nil
}

func (mst *MerkleSearchTree) updateEntry(ctx context.Context, ix int, entry NodeEntry) (*MerkleSearchTree, error) {
	entries, err := mst.getEntries(ctx)
	if err != nil {
		return nil, err
	}

	nents := make([]NodeEntry, len(entries))
	copy(nents, entries[:ix])
	nents[ix] = entry
	copy(nents[ix+1:], entries[ix+1:])

	return mst.newTree(nents), nil
}

func (mst *MerkleSearchTree) replaceWithSplit(ctx context.Context, ix int, left *MerkleSearchTree, nl NodeEntry, right *MerkleSearchTree) (*MerkleSearchTree, error) {
	entries, err := mst.getEntries(ctx)
	if err != nil {
		return nil, err
	}
	var update []NodeEntry
	update = append(update, entries[:ix]...)

	if left != nil {
		update = append(update, NodeEntry{
			Kind: EntryTree,
			Tree: left,
		})
	}

	update = append(update, nl)

	if right != nil {
		update = append(update, NodeEntry{
			Kind: EntryTree,
			Tree: right,
		})
	}

	update = append(update, entries[ix:]...)

	return mst.newTree(update), nil
}

func (mst *MerkleSearchTree) splitAround(ctx context.Context, key []byte) (*MerkleSearchTree, *MerkleSearchTree, error) {
	index, err := mst.findGtOrEqualLeafIndex(ctx, key)
	if err != nil {
		return nil, nil, err
	}

	entries, err := mst.getEntries(ctx)
	if err != nil {
		return nil, nil, err
	}

	leftData := entries[:index]
	rightData := entries[index:]
	left := mst.newTree(leftData)
	right := mst.newTree(rightData)

	if len(leftData) > 0 && leftData[len(leftData)-1].isTree() {
		lastInLeft := leftData[len(leftData)-1]

		nleft, err := left.removeEntry(ctx, len(leftData)-1)
		if err != nil {
			return nil, nil, err
		}
		left = nleft

		subl, subr, err := lastInLeft.Tree.splitAround(ctx, key)
		if err != nil {
			return nil, nil, err
		}

		left, err = left.append(ctx, treeEntry(subl))
		if err != nil {
			return nil, nil, err
		}

		right, err = right.prepend(ctx, treeEntry(subr))
		if err != nil {
			return nil, nil, err
		}

	}

	if left.entryCount() == 0 {
		left = nil
	}
	if right.entryCount() == 0 {
		right = nil
	}
	return left, right, nil
}

func (mst *MerkleSearchTree) entryCount() int {
	entries, err := mst.getEntries(context.TODO())
	if err != nil {
		panic(err)
	}

	return len(entries)
}

func (mst *MerkleSearchTree) append(ctx context.Context, ent NodeEntry) (*MerkleSearchTree, error) {
	entries, err := mst.getEntries(ctx)
	if err != nil {
		return nil, err
	}

	nents := make([]NodeEntry, len(entries)+1)
	copy(nents, entries)
	nents[len(nents)-1] = ent

	return mst.newTree(nents), nil
}

func (mst *MerkleSearchTree) prepend(ctx context.Context, ent NodeEntry) (*MerkleSearchTree, error) {
	entries, err := mst.getEntries(ctx)
	if err != nil {
		return nil, err
	}

	nents := make([]NodeEntry, len(entries)+1)
	copy(nents[1:], entries)
	nents[0] = ent

	return mst.newTree(nents), nil
}

func (mst *MerkleSearchTree) removeEntry(ctx context.Context, ix int) (*MerkleSearchTree, error) {
	entries, err := mst.getEntries(ctx)
	if err != nil {
		return nil, err
	}

	nents := make([]NodeEntry, len(entries)-1)
	copy(nents, entries[:ix])
	copy(nents[ix:], entries[ix+1:])
	return mst.newTree(nents), nil
}

func (mst *MerkleSearchTree) spliceIn(ctx context.Context, entry NodeEntry, ix int) (*MerkleSearchTree, error) {
	entries, err := mst.getEntries(ctx)
	if err != nil {
		return nil, err
	}

	nents := make([]NodeEntry, len(entries)+1)
	copy(nents, entries[:ix])
	nents[ix] = entry
	copy(nents[ix+1:], entries[ix:])

	return mst.newTree(nents), nil
}

func (mst *MerkleSearchTree) atIndex(ix int) (NodeEntry, error) {
	entries, err := mst.getEntries(context.TODO())
	if err != nil {
		return NodeEntry{}, err
	}

	if ix < 0 || ix >= len(entries) {
		return NodeEntry{}, nil
	}

	return entries[ix], nil
}

// this smells inefficient
func (mst *MerkleSearchTree) findGtOrEqualLeafIndex(ctx context.Context, key []byte) (int, error) {
	entries, err := mst.getEntries(ctx)
	if err != nil {
		return -1, err
	}

	for i, e := range entries {
		if e.isLeaf() && bytes.Compare(e.Key, key) > 0 {
			return i, nil
		}
	}

	return len(entries), nil
}

func (mst *MerkleSearchTree) getEntries(ctx context.Context) ([]NodeEntry, error) {
	if mst.entries != nil {
		return mst.entries, nil
	}

	if mst.pointer != cid.Undef {
		var nd NodeData
		if err := mst.cst.Get(ctx, mst.pointer, &nd); err != nil {
			return nil, err
		}

		if len(nd.E) == 0 {
			// maybe this should be an error? idk
			return nil, nil
		}

		firstLeaf := nd.E[0]

		layer := leadingZerosOnHash(firstLeaf.K, mst.fanout)

		entries, err := deserializeNodeData(ctx, mst.cst, &nd, layer, mst.fanout)
		if err != nil {
			return nil, err
		}

		mst.entries = entries
		return entries, nil
	}

	return nil, fmt.Errorf("no entries or cid provided")
}

func (mst *MerkleSearchTree) getPointer(ctx context.Context) (cid.Cid, error) {
	if mst.validPtr {
		return mst.pointer, nil
	}

	entries, err := mst.getEntries(ctx)
	if err != nil {
		return cid.Undef, err
	}

	var outdated []*MerkleSearchTree
	for _, e := range entries {
		if e.Kind == EntryTree && !e.Tree.validPtr {
			outdated = append(outdated, e.Tree)
		}
	}

	if len(outdated) > 0 {
		// this block feels... off
		for _, o := range outdated {
			_, err := o.getPointer(ctx)
			if err != nil {
				return cid.Undef, err
			}
		}
		ne, err := mst.getEntries(ctx)
		if err != nil {
			return cid.Undef, err
		}
		entries = ne
	}

	ptr, err := cidForEntries(ctx, entries, mst.cst)
	if err != nil {
		return cid.Undef, err
	}

	mst.pointer = ptr
	mst.validPtr = true

	return mst.pointer, nil
}

func cidForEntries(ctx context.Context, entries []NodeEntry, cst cbor.IpldStore) (cid.Cid, error) {
	nd, err := serializeNodeData(entries)
	if err != nil {
		return cid.Undef, err
	}

	return cst.Put(ctx, nd)
}

func serializeNodeData(entries []NodeEntry) (*NodeData, error) {
	var data NodeData

	i := 0
	if len(entries) > 0 && entries[0].isTree() {
		i++

		if !entries[0].Tree.validPtr {
			panic("invalid pointer in entries list tree")
		}

		ptr := entries[0].Tree.pointer
		data.L = &ptr
	}

	var lastKey []byte
	for i < len(entries) {
		leaf := entries[i]

		if !leaf.isLeaf() {
			return nil, fmt.Errorf("Not a valid node: two subtrees next to eachother")
		}
		i++

		var subtree *cid.Cid

		if i < len(entries) {
			next := entries[i]

			if next.isTree() {
				subtree = &next.Tree.pointer
				i++
			}

			prefixLen := countPrefixLen(lastKey, leaf.Key)
			data.E = append(data.E, TreeEntry{
				P: int64(prefixLen),
				K: leaf.Key[prefixLen:],
				V: leaf.Val,
				T: subtree,
			})

		}

		lastKey = leaf.Key
	}

	return &data, nil
}

func countPrefixLen(a, b []byte) int {
	// this is probably wrong
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return i
		}
	}

	return len(a)
}

func deserializeNodeData(ctx context.Context, cst cbor.IpldStore, nd *NodeData, layer int, fanout int) ([]NodeEntry, error) {
	var entries []NodeEntry
	if nd.L != nil {
		entries = append(entries, NodeEntry{
			Kind: EntryTree,
			Tree: NewMST(cst, fanout, *nd.L, nil, layer-1),
		})
	}

	var lastKey []byte
	for _, e := range nd.E {
		key := make([]byte, int(e.P)+len(e.K))
		copy(key, lastKey[:e.P])
		copy(key[e.P:], e.K)

		entries = append(entries, NodeEntry{
			Kind: EntryLeaf,
			Key:  key,
			Val:  e.V,
		})

		lastKey = key
		if e.T != nil {
			entries = append(entries, NodeEntry{
				Kind: EntryTree,
				Tree: NewMST(cst, fanout, *e.T, nil, layer-1),
			})
		}
	}

	return entries, nil
}

func layerForEntries(entries []NodeEntry, fanout int) int {
	var firstLeaf NodeEntry
	for _, e := range entries {
		if e.isLeaf() {
			firstLeaf = e
			break
		}
	}

	if firstLeaf.Kind == EntryUndefined {
		return -1
	}

	return leadingZerosOnHash(firstLeaf.Key, fanout)

}

func (mst *MerkleSearchTree) getLayer(ctx context.Context) (int, error) {
	if mst.layer != -1 {
		return mst.layer, nil
	}

	entries, err := mst.getEntries(ctx)
	if err != nil {
		return -1, err
	}

	mst.layer = layerForEntries(entries, mst.fanout)
	return mst.layer, nil
}

func leadingZerosOnHash(k []byte, fanout int) int {
	hv := sha256.Sum256(k)

	var total int
	for i := 0; i < len(hv); i++ {
		n := bits.LeadingZeros8(hv[i])
		total += n
		if n != 8 {
			break
		}
	}
	return total / fanout
}
