package lru

import (
	"errors"
	"fmt"
	"strings"
)

// LRU implements simple LRU cache
type LRU struct {
	list *linkedList
	dict map[int]*node
	cap  int
}

type node struct {
	prev, next *node
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

// NewLRU constructs a new LRU of the given capacity
func NewLRU(capacity int) *LRU {
	return &LRU{
		list: initLL(),
		dict: make(map[int]*node),
		cap:  capacity,
	}
}

// ErrNotFound is not found error
var ErrNotFound = errors.New("not found")

// Get ...
func (l *LRU) Get(k int) (int, error) {
	if n, ok := l.dict[k]; ok {
		l.list.move(n)
		return n.entry.val, nil
	}

	return 0, ErrNotFound
}

// Put ...
func (l *LRU) Put(k, v int) error {
	// existing entry
	if n, ok := l.dict[k]; ok {
		l.dict[k].entry.val = v
		l.list.move(n)
		return nil
	}

	// new entry
	n := &node{entry: &entry{key: k, val: v}}
	l.dict[k] = n
	l.list.addNode(n)
	if l.list.len > l.cap {
		// delete oldest one
		key := l.list.removeOldest()
		delete(l.dict, key)
	}

	return nil
}

func (l *LRU) String() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "%v", l.list.head.next)
	if l.list.head.next != l.list.tail {
		n := l.list.head.next
		for {
			if n.next == l.list.tail {
				break
			}
			fmt.Fprintf(&sb, " => %v", n.next)
			n = n.next
		}
	}

	return sb.String()
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

// move moves the node to head of Linked List
func (ll *linkedList) move(n *node) {
	ll.delNode(n)
	ll.addNode(n)
}

// removeOldest removes the oldest node from Linked List
func (ll *linkedList) removeOldest() int {
	key := ll.tail.prev.entry.key
	ll.delNode(ll.tail.prev)
	return key
}

func (n *node) String() string {
	return fmt.Sprintf("(%d, %d)", n.entry.key, n.entry.val)
}
