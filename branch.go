package oak

import (
	"sync"
)

type branches map[string]*Branch

type Branch struct {
	extant   bool
	Value    interface{}
	Branches branches
	lock     sync.RWMutex
}

func NewBranch() *Branch {
	return &Branch{
		Branches: branches{},
	}
}

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
