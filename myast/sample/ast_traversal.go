package main

import "fmt"

type Example1 struct {
	// Foo Comments
	Foo string `json:"foo"`
}

type Example2 struct {
	// Aoo Comments
	Aoo int `json:"aoo"`
}

// print Hello World
func PrintHello(){
	fmt.Println("Hello World")
}

func Fn1() int32 {
	return int32(0) + 3
}

func Fn2() int32 {
	a := Fn1() + 3
	return a
}