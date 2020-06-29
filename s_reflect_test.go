package mygo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"syscall"
	"testing"
	"unsafe"
)

//https://www.jianshu.com/p/2904efc7f1a8
//https://studygolang.com/articles/22991?fr=sidebar


type SA1 struct {
	A int
	B int
}

func (_ SA1) FnSA1() int{
	return 3
}

type User struct {
	Id   int
	Name string
	Age  int
}

func (u User) ReflectCallFuncHasArgs(name string, age int) {
	fmt.Println("ReflectCallFuncHasArgs name: ", name, ", age:", age, "and origal User.Name:", u.Name)
}

func (u User) ReflectCallFuncNoArgs() {
	fmt.Println("ReflectCallFuncNoArgs")
}

// 通过接口来获取任意参数，然后一一揭晓
func DoFiledAndMethod(input interface{}) {

	getType := reflect.TypeOf(input)
	fmt.Println("get Type is :", getType.Name())

	getValue := reflect.ValueOf(input)
	fmt.Println("get all Fields is:", getValue)

	// 获取方法字段
	// 1. 先获取interface的reflect.Type，然后通过NumField进行遍历
	// 2. 再通过reflect.Type的Field获取其Field
	// 3. 最后通过Field的Interface()得到对应的value
	for i := 0; i < getType.NumField(); i++ {
		field := getType.Field(i)
		value := getValue.Field(i).Interface()
		fmt.Printf("%s: %v = %v\n", field.Name, field.Type, value)
	}

	// 获取方法
	// 1. 先获取interface的reflect.Type，然后通过.NumMethod进行遍历

	for i := 0; i < getType.NumMethod(); i++ {
		m := getType.Method(i)
		fmt.Printf("%s: %v\n", m.Name, m.Type)
	}
}

func TestReflect(t *testing.T) {
	var num float64 = 1.2345

	pointer := reflect.ValueOf(&num)
	value := reflect.ValueOf(num)

	// 可以理解为“强制转换”，但是需要注意的时候，转换的时候，如果转换的类型不完全符合，则直接panic
	// Golang 对类型要求非常严格，类型一定要完全符合
	// 如下两个，一个是*float64，一个是float64，如果弄混，则会panic
	convertPointer := pointer.Interface().(*float64)
	convertValue := value.Interface().(float64)

	fmt.Println(convertPointer)
	fmt.Println(convertValue)
	sa1 := SA1{
		1,
		2,
	}

	DoFiledAndMethod(sa1)

}

func TestSetRealValue(t *testing.T) {
	var num float64 = 1.2345
	fmt.Println("old value of pointer:", num)

	// 通过reflect.ValueOf获取num中的reflect.Value，注意，参数必须是指针才能修改其值
	pointer := reflect.ValueOf(&num)
	newValue := pointer.Elem()

	fmt.Println("type of pointer:", newValue.Type())
	fmt.Println("settability of pointer:", newValue.CanSet())

	// 重新赋值
	newValue.SetFloat(77)
	fmt.Println("new value of pointer:", num)

	////////////////////
	// 如果reflect.ValueOf的参数不是指针，会如何？
	pointer = reflect.ValueOf(num)
	//newValue = pointer.Elem() // 如果非指针，这里直接panic，“panic: reflect: call of reflect.Value.Elem on float64 Value”
}

func TestReflectCall(t *testing.T) {
	user := User{1, "Allen.Wu", 25}

	// 1. 要通过反射来调用起对应的方法，必须要先通过reflect.ValueOf(interface)来获取到reflect.Value，得到“反射类型对象”后才能做下一步处理
	getValue := reflect.ValueOf(user)

	// 一定要指定参数为正确的方法名
	// 2. 先看看带有参数的调用方法
	methodValue := getValue.MethodByName("ReflectCallFuncHasArgs")
	args := []reflect.Value{reflect.ValueOf("wudebao"), reflect.ValueOf(30)}
	methodValue.Call(args)

	// 一定要指定参数为正确的方法名
	// 3. 再看看无参数的调用方法
	methodValue = getValue.MethodByName("ReflectCallFuncNoArgs")
	args = make([]reflect.Value, 0)
	methodValue.Call(args)
}

func fn(a int) int {
	return a+3
}

func replace(target, double uintptr) []byte {
	code := buildJmpDirective(double)
	bytes := entryAddress(target, len(code))
	original := make([]byte, len(bytes))
	copy(original, bytes)
	//modifyBinary(target, code)
	return original
}

func entryAddress(p uintptr, l int) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{Data: p, Len: l, Cap: l}))
}

//func modifyBinary(target uintptr, bytes []byte) {
//	function := entryAddress(target, len(bytes))
//
//	page := entryAddress(pageStart(target), syscall.Getpagesize())
//	err := syscall.Mprotect(page, syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC)
//	if err != nil {
//		panic(err)
//	}
//	copy(function, bytes)
//
//	err = syscall.Mprotect(page, syscall.PROT_READ|syscall.PROT_EXEC)
//	if err != nil {
//		panic(err)
//	}
//}

func pageStart(ptr uintptr) uintptr {
	return ptr & ^(uintptr(syscall.Getpagesize() - 1))
}



func buildJmpDirective(double uintptr) []byte {
	d0 := byte(double)
	d1 := byte(double >> 8)
	d2 := byte(double >> 16)
	d3 := byte(double >> 24)
	d4 := byte(double >> 32)
	d5 := byte(double >> 40)
	d6 := byte(double >> 48)
	d7 := byte(double >> 56)

	return []byte{
		0x48, 0xBA, d0, d1, d2, d3, d4, d5, d6, d7, // MOV rdx, double
		0xFF, 0x22,     // JMP [rdx]
	}
}

type funcValue struct {
	_ uintptr
	p unsafe.Pointer
}

func getPointer(v reflect.Value) unsafe.Pointer {
	return (*funcValue)(unsafe.Pointer(&v)).p
}

func TestReplaceFunc(t *testing.T) {
	var mockfn = func(a int) int { return a+14}
	srcf := reflect.ValueOf(fn)
	tof := reflect.ValueOf(mockfn)
	assert.Equal(t, srcf.Kind(), reflect.Func)
	assert.Equal(t, tof.Kind(), reflect.Func)
	assert.Equal(t, srcf.Type(), tof.Type())

	srcfPointer := *(*uintptr)(getPointer(srcf))
	tofPointer := uintptr(getPointer(tof))

	_ = replace(srcfPointer, tofPointer)
	assert.Equal(t, 15, fn(1))

}

type STM struct {

}

func (_ *STM) fn(a int) int {
	return a+10
}

func TestReplaceMethod(t *testing.T) {
	mockTarget := reflect.TypeOf(&STM{})
	methodName := "Fn"

	m, ok := mockTarget.MethodByName(methodName)
	if !ok {
		panic("retrieve method by name failed")
	}
	fmt.Println(m)
}