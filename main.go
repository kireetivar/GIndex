package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {

	var mut = &sync.Mutex{}
	var wg sync.WaitGroup

	allFiles := TraverseDir("test", 1)
	srcMap := make(map[string]string)
	index := make(map[string][]string)

	wg.Add(len(allFiles))
	for _, file := range allFiles {
		go createSrcMap(file, srcMap,  &wg, mut)
	}
	wg.Wait()

	wg.Add(len(srcMap))
	for fileName, srcCode := range srcMap {
		go createIndex(fileName, srcCode, index, &wg, mut)
	}

	wg.Wait()
	createJsonFile(index)
}

func TraverseDir(s string, depth int) []string {
	dir, err := os.ReadDir(s)
	// TODO handle error properly
	if err != nil {
		fmt.Println("Error reading directory/file:", err)
		os.Exit(1)
	}
	var ds []string

	for _, v := range dir {
		fullpath := filepath.Join(s, v.Name())
		if v.IsDir() {
			ds = append(ds, TraverseDir(fullpath, 1)...)
			depth--
			if depth <= 0 {
				return ds
			}
		} else {
			ds = append(ds, fullpath)
		}
	}
	return ds
}

func createSrcMap(file string, srcMap map[string]string, wg *sync.WaitGroup, mut *sync.Mutex) {
	defer wg.Done()
	bs, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}
	contents := string(bs)
	mut.Lock()
	srcMap[file] = contents
	mut.Unlock()
}

func createIndex(fileName string, srcCode string, index map[string][]string, wg *sync.WaitGroup, mut *sync.Mutex) {
	defer wg.Done()
	keywords := strings.Fields(srcCode)
	for _, word := range keywords {
		mut.Lock()
		value, ok := index[word]
		if ok {
			index[word] = append(value, fileName)
		} else {
			index[word] = []string{fileName}
		}
		mut.Unlock()
	}
}

func createJsonFile(index map[string][]string) {
	jsonbytes ,err := json.MarshalIndent(index, "", "  ")
	if err!=nil {
		fmt.Println("Error converting map to json: ", err)
		os.Exit(1)
	}

	os.WriteFile("jsonIndex.json",jsonbytes,os.ModeAppend)

}