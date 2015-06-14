package main

import (
	"fmt"
	"go-searcher"
	"path/filepath"
	"io/ioutil"
)


func main() {
	searcher := searcher.New(searcher.SearchOptions{
		Cwd:func(path string) string {
			return filepath.Join(".", path)
		},
	})

	index,err := searcher.BuildIndex("files/sources/s3")

	if err != nil {
		fmt.Println(err)
	}

	err = index.ToFile("files/target/t1/from-t1.json")

	if err != nil {
		fmt.Println(err)
	}

	bytes,err := ioutil.ReadFile("files/target/t1/from-t1.json")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(bytes))
}

