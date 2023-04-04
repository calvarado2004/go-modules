package main

import "toolkit"

// main function
func main() {

	var tools toolkit.Tools

	err := tools.CreateDirIfNotExist("./test-dir")
	if err != nil {
		return
	}

}
