package main

import "toolkit"

func main() {

	toSlug := "Now!!? is the time 123"

	var tools toolkit.Tools

	slug, err := tools.Slugify(toSlug)
	if err != nil {
		panic(err)
	}

	println(slug)
}
