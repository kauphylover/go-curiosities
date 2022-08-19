package main

import "fmt"

func main() {
	test()
	fmt.Println("did we get to this?")
}

func test() {
	panic("dying")
}
