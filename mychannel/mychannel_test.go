package mychannel

import (
	"fmt"
	"testing"
	"time"
)

func pass(left, right chan int){
	left <- 1 + <- right
}


func TestMyChannel(t *testing.T) {
	const n = 50
	leftmost := make(chan int)
	right := leftmost
	left := leftmost

	for i := 0; i< n; i++ {
		right = make(chan int)
		// the chain is constructed from the end
		go pass(left, right) // the first goroutine holds (leftmost, new chan)
		left = right         // the second and following goroutines hold (last right chan, new chan)
	}
	go func(c chan int){ c <- 1}(right)
	fmt.Println("sum:", <- leftmost)
}

func TestForRange(t *testing.T) {
	go func() {
		time.Sleep(1 * time.Hour)
	}()
	c := make(chan int)
	go func() {
		for i := 0; i < 10; i = i + 1 {
			c <- i
			time.Sleep(time.Millisecond*500)
		}
		close(c)
	}()
	for i := range c {
		fmt.Println(i)
	}
	fmt.Println("Finished")
}

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum // send sum to c
}

func TestSum(t *testing.T) {
	s := []int{7, 2, 8, -9, 4, 0}
	c := make(chan int)
	go sum(s[:len(s)/2], c)
	go sum(s[len(s)/2:], c)
	x, y := <-c, <-c // receive from c
	fmt.Println(x, y, x+y)
}

func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

func TestSelect(t *testing.T) {
	c := make(chan int)
	quit := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c)
		}
		quit <- 0
	}()
	fibonacci(c, quit)
}

func TestTimeout(t *testing.T) {
	c1 := make(chan string, 1)
	go func() {
		time.Sleep(time.Second * 2)
		c1 <- "result 1"
	}()
	select {
	case res := <-c1:
		fmt.Println(res)
	case <-time.After(time.Second * 1):
		fmt.Println("timeout 1")
	}
}

func TestTimer(t *testing.T) {
	timer1 := time.NewTimer(time.Second * 2)
	<-timer1.C
	fmt.Println("Timer 1 expired")
	timer2 := time.NewTimer(time.Second)
	go func() {
		<-timer2.C
		fmt.Println("Timer 2 expired")
	}()
	stop2 := timer2.Stop()
	if stop2 {
		fmt.Println("Timer 2 stopped")
	}
}

func TestTicker(t *testing.T) {
	ticker := time.NewTicker(time.Millisecond * 500)
	quit := make(chan int)
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
		}
		quit <- 0
	}()
	<- quit
}

func TestClose(t *testing.T) {
	go func() {
		time.Sleep(time.Hour)
	}()
	c := make(chan int, 10)
	c <- 1
	c <- 2
	close(c)
	//c <- 3   //panic  send on closed channel:
}

func TestClose2(t *testing.T) {
	c := make(chan int, 10)
	c <- 1
	c <- 2
	close(c)
	fmt.Println(<-c) //1
	fmt.Println(<-c) //2
	fmt.Println(<-c) //0
	fmt.Println(<-c) //0
}

func TestClose3(t *testing.T) {
	c := make(chan int, 10)
	c <- 1
	c <- 2
	close(c)
	for i := range c {
		fmt.Println(i)  //1,2
	}
}

func worker(done chan bool) {
	time.Sleep(time.Second)
	// 通知任务已完成
	done <- true
}

func TestSync(t *testing.T) {
	done := make(chan bool, 1)
	go worker(done)
	// 等待任务完成
	<-done
}