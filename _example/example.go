package main

import (
	"fmt"
	"github.com/mattn/go-plist"
	"github.com/mattn/go-scan"
	"log"
	"os"
)

func main() {
	f, err := os.Open("example.plist")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	v, err := plist.Read(f)
	if err != nil {
		log.Fatal(err)
	}

	tree := v.(plist.Dict)
	for _, t := range tree["Tracks"].(plist.Dict) {
		if item, ok := t.(plist.Dict); ok {
			fmt.Println(item["Name"])
		}
	}

	var name string
	err = scan.ScanTree(tree, `/Tracks/480/Name`, &name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)
}
