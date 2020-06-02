package main

import (
	"github.com/agiledragon/gomonkey"
)

// Foo 结构体
type Foo struct {
	i int
}

// Bar 接口
type Bar interface {
	Do() error
}

func Fnf1(a int) int {
	return a+3
}

func Fnf2(a int) int {
	return Fnf1(a+3)
}

// main方法
func main() {
	var mocks = func() []*gomonkey.Patches {

		patches := []*gomonkey.Patches{}

		patch1 := gomonkey.ApplyFunc(Fnf1, func(a int) int {
			return 3
		})
		patches = append(patches, patch1)

		patch2 := gomonkey.ApplyFunc(Fnf2, func(a int) int {
			return 3
		})
		patches = append(patches, patch2)

		return patches
	}
	patches := mocks()
	for _, p := range patches {
		p.Reset()
	}
}