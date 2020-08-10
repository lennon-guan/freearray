package freearray

import (
	"fmt"
	"sync"
)

type node struct {
	data           interface{}
	prev, next, at int32
	flags          uint8
}

const (
	flagFree = uint8(0)
	flagBusy = uint8(1)
)

type FreeArray struct {
	lock      sync.Mutex
	items     []node
	allocated int32
	freehead  int32
	busyhead  int32
}

//New Create a new freearray with given capacity
func New(capacity int32) *FreeArray {
	items := make([]node, capacity)
	for i := range items {
		items[i] = node{
			next:  int32(i + 1),
			at:    int32(i),
			prev:  int32(i - 1),
			data:  nil,
			flags: flagFree,
		}
	}
	items[0].prev = -1
	items[capacity-1].next = -1
	return &FreeArray{
		items:    items,
		freehead: 0,
		busyhead: -1,
	}
}

//Alloc Alloc a free node, and returns the index of the node. If alloc nothing, returns -1
func (l *FreeArray) Alloc() int32 {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.freehead < 0 {
		return -1
	}
	index := l.freehead
	node := &l.items[index]
	if (node.flags & flagBusy) != 0 {
		return -1
	}
	l.removeNode(node, &l.freehead)
	l.insertBefore(node, &l.busyhead)
	node.flags |= flagBusy
	l.allocated++
	return index
}

func (l *FreeArray) Release(index int32) {
	l.lock.Lock()
	defer l.lock.Unlock()
	node := &l.items[index]
	if (node.flags & flagBusy) == 0 {
		return
	}
	l.removeNode(node, &l.busyhead)
	l.insertBefore(node, &l.freehead)
	node.flags &= ^flagBusy
	l.allocated--
}

func (l *FreeArray) Data(index int) interface{} {
	if index < 0 || index >= len(l.items) {
		return nil
	}
	if node := l.items[index]; (node.flags & flagBusy) != 0 {
		return node.data
	}
	return nil
}

func (l *FreeArray) removeNode(node *node, headAt *int32) {
	if node.next >= 0 {
		l.items[node.next].prev = node.prev
	}
	if node.prev >= 0 {
		l.items[node.prev].next = node.next
	}
	if *headAt == node.at {
		*headAt = node.next
	}
}

func (l *FreeArray) insertBefore(node *node, headAt *int32) {
	if *headAt < 0 {
		node.next = -1
		node.prev = -1
	} else {
		head := &l.items[*headAt]
		head.prev = node.at
		node.next = *headAt
		node.prev = -1
	}
	*headAt = node.at
}

func (l *FreeArray) printState() {
	fmt.Print("ListState --> Free [")
	for i := l.freehead; i >= 0; i = l.items[i].next {
		fmt.Printf(" %d", i)
	}
	fmt.Print("], ")
	fmt.Print("Busy [")
	for i := l.busyhead; i >= 0; i = l.items[i].next {
		fmt.Printf(" %d", i)
	}
	fmt.Println("]")
}
