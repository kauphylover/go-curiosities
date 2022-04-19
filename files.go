package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	_, err := ioutil.ReadFile("/Users/amallela/.bashrc1")
	if os.IsNotExist(err) {
		fmt.Println("NOT EXIST")
	} else {
		fmt.Println(err)
	}

}
