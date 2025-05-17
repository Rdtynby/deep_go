package main

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

func Trace(stacks [][]uintptr) []uintptr {
	visited := make(map[uintptr]uintptr)
	var result []uintptr

	for _, stack := range stacks {
		for _, pc := range stack {
			loops := make(map[uintptr]uintptr)

			for pc != 0 {
				_, ok := loops[pc]

				if ok {
					break
				}

				loops[pc] = pc

				_, ok2 := visited[pc]

				if !ok2 {
					visited[pc] = pc
					result = append(result, pc)
				}

				pc = uintptr(*(*int)(unsafe.Pointer(pc)))
			}
		}
	}

	return result
}

func TestTrace(t *testing.T) {
	var heapObjects = []int{
		0x00, 0x00, 0x00, 0x00, 0x00,
	}

	var heapPointer1 = &heapObjects[1]
	var heapPointer2 = &heapObjects[2]
	var heapPointer3 *int = nil
	var heapPointer4 = &heapPointer3

	var loopPointer1 *int = nil
	var loopPointer2 = (*int)(unsafe.Pointer(&loopPointer1))
	var loopPointer3 = (*int)(unsafe.Pointer(&loopPointer2))
	loopPointer1 = (*int)(unsafe.Pointer(&loopPointer3))

	var stacks = [][]uintptr{
		{
			uintptr(unsafe.Pointer(&heapPointer1)), 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[0])),
			0x00, 0x00, 0x00, 0x00,
		},
		{
			uintptr(unsafe.Pointer(&heapPointer2)), 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[1])),
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[2])),
			uintptr(unsafe.Pointer(&heapPointer4)), 0x00, 0x00, 0x00,
		},
		{
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&loopPointer1)),
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, uintptr(unsafe.Pointer(&heapObjects[3])),
		},
	}

	pointers := Trace(stacks)
	expectedPointers := []uintptr{
		uintptr(unsafe.Pointer(&heapPointer1)),
		uintptr(unsafe.Pointer(&heapObjects[1])),
		uintptr(unsafe.Pointer(&heapObjects[0])),
		uintptr(unsafe.Pointer(&heapPointer2)),
		uintptr(unsafe.Pointer(&heapObjects[2])),
		uintptr(unsafe.Pointer(&heapPointer4)),
		uintptr(unsafe.Pointer(&heapPointer3)),
		uintptr(unsafe.Pointer(&loopPointer1)),
		uintptr(unsafe.Pointer(&loopPointer3)),
		uintptr(unsafe.Pointer(&loopPointer2)),
		uintptr(unsafe.Pointer(&heapObjects[3])),
	}

	assert.True(t, reflect.DeepEqual(expectedPointers, pointers))
}
