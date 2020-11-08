// +build ignore

package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

var assets http.FileSystem = http.Dir("web/build")

func main() {
	err := vfsgen.Generate(assets, vfsgen.Options{
		Filename:    "api/assets_vfsgen.go",
		PackageName: "api",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
