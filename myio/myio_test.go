package myio

import (
	"fmt"
	"io"
	"testing"
)

type myWriter struct {
	io.Writer
}

func (this *myWriter) Write(p []byte) (n int, err error) {
	//n, err = this.Writer.Write(p)
	fmt.Printf("xdfdfasdf %s", p)
	return n, err
}

func TestWrite(t *testing.T) {
	my := &myWriter{}

	fmt.Fprintf(my, "2323123%s", "2323")

}