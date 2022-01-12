package lfu

import "errors"

// LFU implements simple LFU cache
type LFU struct {
	listmap map[int]*linkedList
	dict    map[int]*node
	total   int // node of list map
	minFreq int
	cap     int
}

type node struct {
	prev, next *node
	freq       int
	entry      *entry
}

type entry struct {
	key int
	val int
}

type linkedList struct {
	head, tail *node
	len        int
}

// NewLFU constructs a new LFU of the given capacity
func NewLFU(capacity int) *LFU {
	return &LFU{
		listmap: make(map[int]*linkedList),
		dict:    make(map[int]*node),
		total:   0,
		minFreq: 0,
		cap:     capacity,
	}
}

// ErrNotFound is not found error
var ErrNotFound = errors.New("not found")

// Get ...
func (l *LFU) Get(k int) (int, error) {
	if n, ok := l.dict[k]; ok {
		l.update(n)
		return n.entry.val, nil
	}

	return 0, ErrNotFound
}

// Put ...
func (l *LFU) Put(k, v int) error {
	// existing entry
	if n, ok := l.dict[k]; ok {
		l.dict[k].entry.val = v
		l.update(n)
		return nil
	}

	// new entry
	if l.total == l.cap {
		oldest := l.listmap[l.minFreq].tail.prev
		l.listmap[l.minFreq].delNode(oldest)
		delete(l.dict, oldest.entry.key)
		l.total--
	}
	n := &node{entry: &entry{key: k, val: v}, freq: 1}
	l.dict[k] = n
	if _, ok := l.listmap[1]; !ok {
		l.listmap[1] = initLL()
	}
	l.listmap[1].addNode(n)
	l.minFreq = 1
	l.total++

	return nil
}

func (l *LFU) update(n *node) {
	curFreq := n.freq
	l.listmap[curFreq].delNode(n)
	n.freq = curFreq + 1
	if _, ok := l.listmap[curFreq+1]; !ok {
		l.listmap[curFreq+1] = initLL()
	}
	l.listmap[curFreq+1].addNode(n)
	if curFreq == l.minFreq && l.listmap[curFreq].len == 0 {
		l.minFreq = curFreq + 1
	}
}

func initLL() *linkedList {
	ll := &linkedList{
		head: &node{},
		tail: &node{},
	}
	ll.head.next = ll.tail
	ll.tail.prev = ll.head

	return ll
}

// addNode adds new node in front of list
func (ll *linkedList) addNode(n *node) {
	n.prev = ll.head
	n.next = ll.head.next

	ll.head.next = n
	n.next.prev = n

	ll.len++
}

// delNode deletes node from Linked List
func (ll *linkedList) delNode(n *node) {
	n.prev.next = n.next
	n.next.prev = n.prev

	ll.len--
}
