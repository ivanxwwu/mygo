package other1

import (
	"fmt"
)

type TestPointer struct {
	A int32
	b int32

}

func (T *TestPointer) OutPut() {
	fmt.Println("TestPointer OutPut:", T.A, T.b)
}

func privateFunc() {
	fmt.Println("this is private func\n")
}

type PublicStruct struct {
	I int
	b int
}

func (p *PublicStruct) f(a int) {
	fmt.Printf("PublicStruct f() %d, %d\n", p.b, p.I)
}

var private_m = map[int]string {
	1:"a",
}