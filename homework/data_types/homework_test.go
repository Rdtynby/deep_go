package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// go test -v homework_test.go

func ToLittleEndian(number uint32) uint32 {
	return number&0xFF<<24 + number&(0xFF<<8)<<8 + number&(0xFF<<16)>>8 + number&(0xFF<<24)>>24
}

func TestConversion(t *testing.T) {
	tests := map[string]struct {
		number uint32
		result uint32
	}{
		"test case #1": {
			number: 0x00000000,
			result: 0x00000000,
		},
		"test case #2": {
			number: 0xFFFFFFFF,
			result: 0xFFFFFFFF,
		},
		"test case #3": {
			number: 0x00FF00FF,
			result: 0xFF00FF00,
		},
		"test case #4": {
			number: 0x0000FFFF,
			result: 0xFFFF0000,
		},
		"test case #5": {
			number: 0x01020304,
			result: 0x04030201,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndian(test.number)
			assert.Equal(t, test.result, result)
		})
	}
}

type Number interface {
	uint8 | uint16 | uint32 | uint64
}

func ToLittleEndianGeneric[T Number](number T, lenInBytes int) T {
	if lenInBytes == 1 {
		return number
	}

	lenInBytes /= 2
	shift := lenInBytes * 8
	part1 := ((1<<shift - 1) & number << shift) >> shift
	part2 := number >> shift

	return ToLittleEndianGeneric[T](part1, lenInBytes)<<shift + ToLittleEndianGeneric[T](part2, lenInBytes)
}

func TestConversionGeneric(t *testing.T) {
	tests64 := map[string]struct {
		number uint64
		result uint64
	}{
		"test case 64 #1": {
			number: 0x0000000000000000,
			result: 0x0000000000000000,
		},
		"test case 64 #2": {
			number: 0xFFFFFFFFFFFFFFFF,
			result: 0xFFFFFFFFFFFFFFFF,
		},
		"test case 64 #3": {
			number: 0x00FF00FF00FF00FF,
			result: 0xFF00FF00FF00FF00,
		},
		"test case 64 #4": {
			number: 0x0000FFFF0000FFFF,
			result: 0xFFFF0000FFFF0000,
		},
		"test case 64 #5": {
			number: 0x0102030405060708,
			result: 0x0807060504030201,
		},
	}

	for name, test := range tests64 {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndianGeneric(test.number, 8)
			assert.Equal(t, test.result, result)
		})
	}

	tests32 := map[string]struct {
		number uint32
		result uint32
	}{
		"test case 32 #1": {
			number: 0x00000000,
			result: 0x00000000,
		},
		"test case 32 #2": {
			number: 0xFFFFFFFF,
			result: 0xFFFFFFFF,
		},
		"test case 32 #3": {
			number: 0x00FF00FF,
			result: 0xFF00FF00,
		},
		"test case 32 #4": {
			number: 0x0000FFFF,
			result: 0xFFFF0000,
		},
		"test case 32 #5": {
			number: 0x01020304,
			result: 0x04030201,
		},
	}

	for name, test := range tests32 {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndianGeneric(test.number, 4)
			assert.Equal(t, test.result, result)
		})
	}

	tests16 := map[string]struct {
		number uint16
		result uint16
	}{
		"test case 16 #1": {
			number: 0x0000,
			result: 0x0000,
		},
		"test case 16 #2": {
			number: 0xFFFF,
			result: 0xFFFF,
		},
		"test case 16 #3": {
			number: 0x00FF,
			result: 0xFF00,
		},
		"test case 16 #4": {
			number: 0x0102,
			result: 0x0201,
		},
	}

	for name, test := range tests16 {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndianGeneric(test.number, 2)
			assert.Equal(t, test.result, result)
		})
	}

	tests8 := map[string]struct {
		number uint8
		result uint8
	}{
		"test case 8 #1": {
			number: 0x00,
			result: 0x00,
		},
		"test case 8 #2": {
			number: 0xFF,
			result: 0xFF,
		},
	}

	for name, test := range tests8 {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndianGeneric(test.number, 1)
			assert.Equal(t, test.result, result)
		})
	}
}
