package main

import (
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"log"
)

var data = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`

// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.
type T struct {
	A string
	B struct {
		RenamedC int   `yaml:"c"`
		D        []int `yaml:",flow"`
	}
}

var dms = `
dm1:
  url: git1
  default: true
dm2:
  url: fs1
  default: false
`

type Props struct {
	Url     string `yaml:"url"`
	Default bool   `yaml:"default"`
}

func main() {
	e := make(map[string]Props)
	p1 := Props{"url1", true}
	e["dm1"] = p1
	e["dm2"] = p1
	//dms2 := DmMap{e}
	b, err1 := yaml.Marshal(&e)
	if err1 == nil {
		fmt.Printf("Marshaled\n%v\n\n", string(b))
	} else {
		fmt.Printf("#### ERROR: %v\n\n", err1)
	}

	var eu map[string]Props = make(map[string]Props)
	err2 := yaml.Unmarshal([]byte(dms), &eu)
	if err2 != nil {
		fmt.Printf("### ERROR Unmarshaling: %v\n", err2)
		return
	} else {
		fmt.Printf("Success: \n%v\n\n", eu)
		for key, props := range eu {
			fmt.Printf("%v: %v\n", key, props)
		}
	}

	t := T{}

	err := yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t:\n%v\n\n", t)

	d, err := yaml.Marshal(&t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t dump:\n%s\n\n", string(d))

	m := make(map[interface{}]interface{})

	err = yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- m:\n%v\n\n", m)

	d, err = yaml.Marshal(&m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- m dump:\n%s\n\n", string(d))
}
