package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Number interface {
	int8 | int16 | int32 | int64
}

type Node[T Number] struct {
	key   T
	value T
	left  *Node[T]
	right *Node[T]
}

type OrderedMap[T Number] struct {
	size int
	root *Node[T]
}

func NewOrderedMap[T Number]() OrderedMap[T] {
	return OrderedMap[T]{}
}

func InsertNode[T Number](root **Node[T], newNode *Node[T]) {
	if *root == nil {
		*root = newNode
	} else if (*root).key > newNode.key {
		InsertNode(&(*root).left, newNode)
	} else {
		InsertNode(&(*root).right, newNode)
	}
}

func (m *OrderedMap[T]) Insert(key T, value T) {
	m.size++

	InsertNode(&m.root, &Node[T]{key: key, value: value})
}

func (root *Node[T]) FindMinChild() *Node[T] {
	if root == nil {
		return nil
	} else if root.left == nil {
		return root
	} else {
		return root.left.FindMinChild()
	}
}

func EraseNode[T Number](root **Node[T], key T) {
	if *root == nil {
		return
	}

	if (*root).key == key {
		if (*root).left == nil {
			*root = (*root).right
		} else if (*root).right == nil {
			*root = (*root).left
		} else {
			minChild := (*root).right.FindMinChild()
			(*root).value = minChild.value
			(*root).key = minChild.key
			EraseNode(&minChild.right, minChild.key)
		}
	} else if (*root).key > key {
		EraseNode(&(*root).left, key)
	} else {
		EraseNode(&(*root).right, key)
	}
}

func (m *OrderedMap[T]) Erase(key T) {
	m.size--

	EraseNode(&m.root, key)
}

func (root *Node[T]) ContainsNode(key T) bool {
	if root == nil {
		return false
	}

	if (*root).key > key {
		return (*root).left.ContainsNode(key)
	} else if (*root).key < key {
		return (*root).right.ContainsNode(key)
	} else {
		return true
	}
}

func (m *OrderedMap[T]) Contains(key T) bool {
	return m.root.ContainsNode(key)
}

func (m *OrderedMap[T]) Size() int {
	return m.size
}

func (root *Node[T]) ForEachNode(action func(T, T)) {
	if root == nil {
		return
	}

	root.left.ForEachNode(action)
	action(root.key, root.value)
	root.right.ForEachNode(action)
}

func (m *OrderedMap[T]) ForEach(action func(T, T)) {
	m.root.ForEachNode(action)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap[int64]()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int64
	expectedKeys := []int64{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int64) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int64{4, 5, 10, 12}
	data.ForEach(func(key, _ int64) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
