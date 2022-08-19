package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	var asdf = make(map[string]any)
	asdf["asdf"] = 1
	asdf["123"] = "qwer"

	fmt.Println(asdf["12"].(string))

	_, err := ioutil.ReadFile("/Users/amallela/.bashrc1")
	if os.IsNotExist(err) {
		fmt.Println("NOT EXIST")
	} else {
		fmt.Println(err)
	}

}
