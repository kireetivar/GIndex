package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	AllFiles := TraverseDir("test",1)

	index := createIndex(AllFiles)

	fmt.Println(index)
}

func TraverseDir(s string, depth int) []string{
	dir, err := os.ReadDir(s)
	// TODO handle error properly
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}
	var ds []string

	for _, v := range dir {
		fullpath := filepath.Join(s, v.Name())
		if v.IsDir() {
			ds = append(ds, TraverseDir(fullpath, 1)...)
			depth--;
			if depth<=0{
				return ds
			}
		}else{
			ds = append(ds, fullpath)
		}
	}
	return ds
}


func createIndex(allfiles []string) map[string][]string{
	srcMap := make(map[string]string)
	for _,v := range allfiles{
		bs,err := os.ReadFile(v)
		if err != nil {
			fmt.Println("Error reading file:", err)
			os.Exit(1)
		}
		contents := string(bs)
		srcMap[v] = contents
	}

	index := make(map[string][]string)


	for k,v := range srcMap {
		keywords := strings.Fields(v)


		for _, v := range keywords {
			value ,ok := index[v]
			if ok {
				index[v] = append(value, k)
			}else {
				index[v] = []string{k}
			}
		}
	}

	return index
}
