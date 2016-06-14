// Package pair defines pair/tuple types
package pair

//go:generate rxgen -type Int32Pair -name RxInt32Pair .

// Int32Pair is a pair of integers
type Int32Pair struct {
	L int32
	R int32
}
