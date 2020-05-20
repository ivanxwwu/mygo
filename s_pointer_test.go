package mygo

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestPointer1(t *testing.T){
	var x struct {
		a bool
		b int16
		c []int
	}

	// 和 pb := &x.b 等价
	pb := (*int16)(unsafe.Pointer(
		uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)))
	*pb = 42
	fmt.Println(x.b)
}

func TestPointer2(t *testing.T) {
	var float64bits = func (f float64) uint64 { return *(*uint64)(unsafe.Pointer(&f)) }
	fmt.Printf("%#016x\n", float64bits(1.0)) // "0x3ff0000000000000"
}

func TestPointer3(t *testing.T) {
	var x struct {
		a bool
		b int16
		c []int
	}

	// 和 pb := &x.b 等价
	pb := (*int16)(unsafe.Pointer(
		uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)))
	*pb = 42
	fmt.Println(x.b) // "42"
}