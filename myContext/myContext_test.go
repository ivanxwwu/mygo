package myContext

import (
	"context"
	"log"
	"os"
	"testing"
	"time"
)

var logg *log.Logger

func someHandler() {
	ctx, cancel := context.WithCancel(context.Background())
	go doStuff(ctx)

	//10秒后取消doStuff
	time.Sleep(10 * time.Second)
	cancel()

}

func timeoutHandler() {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	// go doTimeOutStuff(ctx)
	go doStuff(ctx)

	time.Sleep(10 * time.Second)

	cancel()
}

func doTimeOutStuff(ctx context.Context) {
	for {
		time.Sleep(1 * time.Second)

		if deadline, ok := ctx.Deadline(); ok { //设置了deadl
			logg.Printf("deadline set")
			if time.Now().After(deadline) {
				logg.Printf(ctx.Err().Error())
				return
			}
		}
		select {
		case <-ctx.Done():
			logg.Printf("done")
			return
		default:
			logg.Printf("work")
		}
	}
}

//每1秒work一下，同时会判断ctx是否被取消了，如果是就退出
func doStuff(ctx context.Context) {
	for {
		time.Sleep(1 * time.Second)
		select {
		case <-ctx.Done():
			logg.Printf("done")
			return
		default:
			logg.Printf("work")
		}
	}
}

func TestSomeHandler(t *testing.T) {
	logg = log.New(os.Stdout, "", log.Ltime)
	someHandler()
	logg.Printf("down")
}

func TestTimeoutHandler(t *testing.T) {
	logg = log.New(os.Stdout, "", log.Ltime)
	timeoutHandler()
	logg.Printf("down")
}

func timeoutHandler2() {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	go doTimeOutStuff(ctx)

	time.Sleep(10 * time.Second)

	cancel()

}

func TestTimeoutHandler2(t *testing.T) {
	logg = log.New(os.Stdout, "", log.Ltime)
	timeoutHandler2()
	logg.Printf("down")
}

