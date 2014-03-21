package oak

import (
	"sync"
)

type branches map[string]*Branch

// Branch represents a node in the tree, and can be either a root node
// or an internal one.
type Branch struct {
	extant   bool
	Value    interface{}
	Branches branches
	lock     sync.RWMutex
}

// NewBranch returns a pointer to a new Branch struct
func NewBranch() *Branch {
	return &Branch{
		Branches: branches{},
	}
}

// Get accepts a ([]string) and returns the match value, and a boolean to
// indicate whether a match was found. 
func (b *Branch) Get(path []string) (value interface{}, exists bool) {
	if len(path) == 0 {
		return b.ownValue()
	}

	key, nesting := path[0], path[1:]

	b.lock.RLock()
	defer b.lock.RUnlock()

	res, ok := b.Branches[key]
	if !ok {
		return nil, false
	}

	return res.Get(nesting)
}

// Set accepts a path ([]string) and a value (interface{}) and sets the
// value of the relevant node in the tree to that interface.
func (b *Branch) Set(path []string, value interface{}) {
	if len(path) == 0 {
		b.setValue(value)
		return
	}

	key, nesting := path[0], path[1:]

	b.lock.RLock()
	res, ok := b.Branches[key]
	b.lock.RUnlock()

	if !ok {
		res = NewBranch()
		b.lock.Lock()
		b.Branches[key] = res
		b.lock.Unlock()
	}

	b.lock.Lock()
	res.Set(nesting, value)
	b.lock.Unlock()
}

// Replace accepts a path ([]string), and replaces a found node with the
// provided *Branch, returning true. If it cannot find a suitable node 
// to replace then it will return false. You cannot replace the root node,
// but this will update its values to match the provided Branch.
func (b *Branch) Replace(path []string, branch *Branch) bool {
	if len(path) == 0 {
		b.lock.Lock()
		defer b.lock.Unlock()
		b.Branches = branch.Branches
		b.Value = branch.Value
		return true
	}

	key, nesting := path[0], path[1:]

	b.lock.Lock()
	defer b.lock.Unlock()

	res, ok := b.Branches[key]

	if !ok {
		return false
	}

	return res.Replace(nesting, branch)
}

// Remove accepts a path, and will remove that node and all nested nodes
// below it from the tree, returning a bool to indicate success. You 
// cannot remove the root node.
func (b *Branch) Remove(path []string) bool {
	switch len(path) {
	case 0:
		return false
	case 1:
		key := path[0]
		b.lock.RLock()
		if _, ok := b.Branches[key]; ok {
			b.lock.RUnlock()
			b.lock.Lock()
			defer b.lock.Unlock()
			delete(b.Branches, key)
			return true
		}
		b.lock.RUnlock()
		return false
	default:
		key, nesting := path[0], path[1:]
		b.lock.RLock()
		if res, ok := b.Branches[key]; ok {
			b.lock.RUnlock()
			b.lock.Lock()
			defer b.lock.Unlock()
			return res.Remove(nesting)
		}
		b.lock.RUnlock()
		return false
	}
}

func (b *Branch) ownValue() (value interface{}, exists bool) {
	if b.extant {
		return b.Value, true
	}
	return nil, false
}

func (b *Branch) setValue(value interface{}) {
	b.extant = true
	b.Value = value
}

func (b *Branch) removeValue(value interface{}) (outcome bool) {
	outcome = b.extant
	b.extant = false
	b.Value = nil
	return outcome
}
