package main

import (
	"fmt"
	"github.com/calvarado2004/go-modules/toolkit"
)

func main() {
	var tools toolkit.Tools

	s := tools.RandomString(10)
	fmt.Printf("Random string %s", s)

}
