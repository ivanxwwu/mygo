package myerror

import (
	"fmt"
	"github.com/pkg/errors"
	"testing"
)

func A()  error {
	return B()
}

func B() error {
	return errors.Wrap(errors.New("fuck err"), "open failed")
}

func TestMyError(t *testing.T) {
	err := B()
	fmt.Printf("original errors:%T %v\n", errors.Cause(err), errors.Cause(err))
	fmt.Printf("%+v\n", err)
}