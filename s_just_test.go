package mygo

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

type Slice []int

func (A Slice)Append(value int) {
	A1 := append(A, value)

	sh:=(*reflect.SliceHeader)(unsafe.Pointer(&A))
	fmt.Printf("A Data:%d,Len:%d,Cap:%d\n",sh.Data,sh.Len,sh.Cap)

	sh1:=(*reflect.SliceHeader)(unsafe.Pointer(&A1))
	fmt.Printf("A1 Data:%d,Len:%d,Cap:%d\n",sh1.Data,sh1.Len,sh1.Cap)
}


func TestUnknown(t *testing.T) {
	mSlice := make(Slice, 10, 20)
	mSlice.Append(5)
	mSlice[1] = 1
	fmt.Println(mSlice)
}