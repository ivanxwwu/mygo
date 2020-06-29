package main

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
	"os"
	"runtime/trace"
)

//https://www.cnblogs.com/Dr-wei/p/11742414.html
func main() {
	//go func() {
	//	i := 0
	//	for {
	//		i++
	//		s1 := fmt.Sprintf("data1:%s len:%d\n", data.Add("https://github.com/EDDYCJY"), data.Len1)
	//		s2 := fmt.Sprintf("data2:%s len:%d\n", data.Add2("https://github.com/22233"), data.Len2)
	//		if i % 1000==0 {
	//			log.Println(s1)
	//			log.Println(s2)
	//		}
	//	}
	//}()
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()
	doWork := func(strings <-chan string) <-chan interface {} {
		completed := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(completed)
			for s := range strings{
				//做一些有趣的事
				fmt.Println(s)
			}
			fmt.Println(11)
		}()
		return completed
	}
	for i:=0; i<20; i++ {
		doWork(nil)
	}

	mux := http.NewServeMux()
	mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

	http.ListenAndServe("0.0.0.0:6060", mux)
}


