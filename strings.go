package main

import "fmt"
import "github.com/coreos/go-semver/semver"

func main() {
	v1 := semver.New("0.1.11")
	v2 := semver.New("0.1.12")
	v1a := semver.New("0.1.0-a")
	v1b := semver.New("0.1.0-b")

	fmt.Println(v1.LessThan(*v2))
	fmt.Println(v1.Compare(*v1a))
	fmt.Println(v1a.LessThan(*v1b))
}
