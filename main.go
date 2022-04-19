package main

import (
	"fmt"
)

func main() {
	fmt.Println("hello world")

	fmt.Printf("Version=%s\n", Get().Version)
	fmt.Printf("GitCommit=%s\n", Get().GitCommit)
	fmt.Printf("Go Version=%s\n", Get().GoVersion)
}
