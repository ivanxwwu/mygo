package access_private

import (
	"fmt"
	"github.com/ivanxwwu/mygo/access_private/other1"
	_ "github.com/ivanxwwu/mygo/access_private/other1"
	"testing"
	"unsafe"
	_ "unsafe"
)

func TestAccessLower(t *testing.T) {
	x := other1.TestPointer{}
	fmt.Printf("size:%d\n", unsafe.Sizeof(x))
	pb := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&x))+4))
	*pb = 2
	x.OutPut()
}

type SizeOfA struct {
	A int
}

type SizeOfC struct {
	A byte
	C int32
}

type SizeOfF struct {
	A byte
	C int16
	B int64
	D int16
}

func (this *SizeOfF) Fn()  int {
	return 0
}

func TestMemoryOffset(t *testing.T) {
	fmt.Printf("sizeofc:%d, alignof:%d\n", unsafe.Sizeof(SizeOfC{0,0}), unsafe.Alignof(SizeOfC{0,0}))
	fmt.Printf("sizeoff:%d, alignoff:%d\n", unsafe.Sizeof(SizeOfF{0,0,0,0}), unsafe.Alignof(SizeOfF{0,0,0,0}))
	fmt.Printf("offsetoffa:%d,offsetoffb:%d,offsetoffc:%d,offsetoffd:%d,", unsafe.Offsetof(SizeOfF{}.A),unsafe.Offsetof(SizeOfF{}.C),
		unsafe.Offsetof(SizeOfF{}.B),unsafe.Offsetof(SizeOfF{}.D))
	//fmt.Printf("offsetoff fn:%d", unsafe.Offsetof((&SizeOfF{}).Fn))
}

//go:linkname private_func github.com/ivanxwwu/mygo/access_private/other1.privateFunc
func private_func()

//go:linkname pulibc_struct_private_func github.com/ivanxwwu/mygo/access_private/other1.(*PublicStruct).f
func pulibc_struct_private_func(p *other1.PublicStruct, a int)

//go:linkname private_member github.com/ivanxwwu/mygo/access_private/other1.private_m
var private_member map[int]string

func TestPrivateFunc(t *testing.T) {
	private_func()
	p := &other1.PublicStruct{I:1}
	pulibc_struct_private_func(p, 100)
	fmt.Printf("private_member:%v\n",  private_member)
}

