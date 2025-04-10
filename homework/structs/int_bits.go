package main

type Number interface {
	uint8 | uint32
}

func getBitsFromInt[T Number](value T, start int, length int) T {
	mask := T(1<<length - 1)
	return value >> start & mask
}

func setBitsToInt[T Number](value T, bits T, start int, length int) T {
	mask := T(1<<length-1) << start
	return value & ^mask | bits<<start
}
