package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func main() {

	f := "./templates/bookinfo-apps.tmpl"
	absPath, err := filepath.Abs(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var tpl bytes.Buffer
	template.Must(template.ParseFiles(absPath)).Execute(&tpl, struct {
		Namespace       string
		DeployReviewsV1 bool
	}{
		Namespace:       "amallela",
		DeployReviewsV1: false,
	})
	fmt.Println(tpl.String())

}
