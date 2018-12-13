package main

import (
	"fmt"

	"github.com/zapkub/gorogoso/test/nested"
)

func main() {
	a := nested.NewA()
	fmt.Printf("Hello world value = %s \n", a.Value)
}
// hi