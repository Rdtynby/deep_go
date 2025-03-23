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

type OrderedMap[T Number] struct {
	key   T
	value T
	size  int
	left  *OrderedMap[T]
	right *OrderedMap[T]
}

func NewOrderedMap[T Number]() OrderedMap[T] {
	return OrderedMap[T]{}
}

func (m *OrderedMap[T]) Insert(key T, value T) {
	m.size++

	if m.key == 0 {
		m.key = key
		m.value = value
	} else {
		newNode := OrderedMap[T]{key: key, value: value}
		parent := (*OrderedMap[T])(nil)

		for m != nil {
			parent = m

			if m.key > key {
				m = m.left
			} else if m.key < key {
				m = m.right
			}
		}

		if parent.key > key {
			parent.left = &newNode
		} else {
			parent.right = &newNode
		}
	}
}

func (m *OrderedMap[T]) Erase(key T) {
	parent := (*OrderedMap[T])(nil)
	m.size--

	for m != nil {
		if m.key == key {
			break
		}

		parent = m

		if m.key > key {
			m = m.left
		} else if m.key < key {
			m = m.right
		}
	}

	if m.right == nil {
		if parent.key < m.key {
			parent.right = m.left
		} else {
			parent.left = m.left
		}
	} else {
		minimum := m.right
		parent = nil

		for minimum.left != nil {
			parent = minimum
			minimum = minimum.left
		}

		if parent != nil {
			parent.left = minimum.right
		} else {
			m.right = minimum.right
		}

		m.key = minimum.key
		m.value = minimum.value
	}
}

func (m *OrderedMap[T]) Contains(key T) bool {
	if m.key == key {
		return true
	}

	for m != nil {
		if m.key > key {
			m = m.left
		} else if m.key < key {
			m = m.right
		} else {
			return true
		}
	}

	return false
}

func (m *OrderedMap[T]) Size() int {
	return m.size
}

func (m *OrderedMap[T]) ForEach(action func(T, T)) {
	if m.left != nil {
		m.left.ForEach(action)
	}

	action(m.key, m.value)

	if m.right != nil {
		m.right.ForEach(action)
	}
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
