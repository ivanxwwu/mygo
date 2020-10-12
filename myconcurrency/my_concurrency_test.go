package myconcurrency

import (
	"fmt"
	"testing"
	"time"
)

func fn1(stop <- chan int) {
	fmt.Println("fn1 start!")
	<-stop
	fmt.Println("fn1 over!")
}

func fn2(stop <- chan int) {
	fmt.Println("fn2 start!")
	<-stop
	fmt.Println("fn2 over!")
}

type Student struct {
	name string
}

func (s *Student) Name() string {
	if s == nil {
		fmt.Println(111111)
		return "3333"
	}
	return "111"
}

func TestStop(t *testing.T) {
	stop := make(chan int)
	go fn1(stop)
	go fn2(stop)
	time.Sleep(time.Second*3)
	//stop <- 1
	close(stop)
}

func TestStudent(t *testing.T) {
	var stu *Student = nil
	fmt.Println(stu.Name())
}