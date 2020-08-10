package freearray

import (
	"fmt"
	"sync"
)

type node struct {
	data       interface{}
	prev, next int32
}

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
			next: int32(i + 1),
			prev: int32(i - 1),
			data: nil,
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
func (l *FreeArray) Alloc(data interface{}) int32 {
	if data == nil {
		panic("data cannot be nil")
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.freehead < 0 {
		return -1
	}
	index := l.freehead
	node := &l.items[index]
	if node.data != nil {
		return -1
	}
	l.removeNode(node, index, &l.freehead)
	l.insertBefore(node, index, &l.busyhead)
	node.data = data
	l.allocated++
	return index
}

func (l *FreeArray) Release(index int32) {
	l.lock.Lock()
	defer l.lock.Unlock()
	node := &l.items[index]
	if node.data == nil {
		return
	}
	l.removeNode(node, index, &l.busyhead)
	l.insertBefore(node, index, &l.freehead)
	node.data = nil
	l.allocated--
}

func (l *FreeArray) Data(index int) interface{} {
	if index < 0 || index >= len(l.items) {
		return nil
	}
	return l.items[index]
}

func (l *FreeArray) removeNode(node *node, index int32, headAt *int32) {
	if node.next >= 0 {
		l.items[node.next].prev = node.prev
	}
	if node.prev >= 0 {
		l.items[node.prev].next = node.next
	}
	if *headAt == index {
		*headAt = node.next
	}
}

func (l *FreeArray) insertBefore(node *node, index int32, headAt *int32) {
	if *headAt < 0 {
		node.next = -1
		node.prev = -1
	} else {
		head := &l.items[*headAt]
		head.prev = index
		node.next = *headAt
		node.prev = -1
	}
	*headAt = index
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
