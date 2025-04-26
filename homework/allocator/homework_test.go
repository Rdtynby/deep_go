package main

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Number interface {
	byte | int | int8 | int16 | int32 | int64
}

func Defragment[T Number](memory []T, pointers []unsafe.Pointer) {
	currentPointer := 0
	changes := make(map[unsafe.Pointer]unsafe.Pointer)

	for index, pointer := range pointers {
		if changes[pointer] == nil {
			changes[pointer] = unsafe.Pointer(&memory[currentPointer])

			memory[currentPointer], *(*T)(pointers[index]) = *(*T)(pointers[index]), memory[currentPointer]

			currentPointer++
		}

		pointers[index] = changes[pointer]
	}
}

func TestDefragmentation(t *testing.T) {
	var fragmentedMemory = []byte{
		0xF0, 0x00, 0x00, 0x00,
		0x00, 0xF1, 0x00, 0x00,
		0x00, 0x00, 0xF2, 0x00,
		0x00, 0x00, 0xF4, 0xF3,
	}

	var fragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[5]),
		unsafe.Pointer(&fragmentedMemory[10]),
		unsafe.Pointer(&fragmentedMemory[15]),
		unsafe.Pointer(&fragmentedMemory[10]),
		unsafe.Pointer(&fragmentedMemory[14]),
	}

	var defragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[1]),
		unsafe.Pointer(&fragmentedMemory[2]),
		unsafe.Pointer(&fragmentedMemory[3]),
		unsafe.Pointer(&fragmentedMemory[2]),
		unsafe.Pointer(&fragmentedMemory[4]),
	}

	var defragmentedMemory = []byte{
		0xF0, 0xF1, 0xF2, 0xF3,
		0xF4, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	Defragment(fragmentedMemory, fragmentedPointers)
	assert.True(t, reflect.DeepEqual(defragmentedMemory, fragmentedMemory))
	assert.True(t, reflect.DeepEqual(defragmentedPointers, fragmentedPointers))
}

func TestDefragmentationInt64(t *testing.T) {
	var fragmentedMemory = []int64{
		0xF0, 0x00, 0x00, 0x00,
		0x00, 0xF1, 0x00, 0x00,
		0x00, 0x00, 0xF2, 0x00,
		0x00, 0x00, 0xF4, 0xF3,
	}

	var fragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[5]),
		unsafe.Pointer(&fragmentedMemory[10]),
		unsafe.Pointer(&fragmentedMemory[15]),
		unsafe.Pointer(&fragmentedMemory[10]),
		unsafe.Pointer(&fragmentedMemory[14]),
	}

	var defragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[1]),
		unsafe.Pointer(&fragmentedMemory[2]),
		unsafe.Pointer(&fragmentedMemory[3]),
		unsafe.Pointer(&fragmentedMemory[2]),
		unsafe.Pointer(&fragmentedMemory[4]),
	}

	var defragmentedMemory = []int64{
		0xF0, 0xF1, 0xF2, 0xF3,
		0xF4, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	Defragment(fragmentedMemory, fragmentedPointers)
	assert.True(t, reflect.DeepEqual(defragmentedMemory, fragmentedMemory))
	assert.True(t, reflect.DeepEqual(defragmentedPointers, fragmentedPointers))
}
