package main

import (
	"fmt"
	"toolkit"
)

func main() {
	var tools toolkit.Tools

	s := tools.RandomString(10)
	fmt.Printf("Random string %s", s)

}
