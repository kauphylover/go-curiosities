package main

import (
	"fmt"
	"go-curiosities/pkg/version"
)

func main() {
	fmt.Println("hello world")

	fmt.Printf("Version=%s\n", version.Get().Version)
	fmt.Printf("GitCommit=%s\n", version.Get().GitCommit)
	fmt.Printf("Go Version=%s\n", version.Get().GoVersion)
}
