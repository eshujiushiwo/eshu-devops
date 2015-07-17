package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var xxx = make(map[string]interface{})

func readfile(filename string) (map[string]interface{}, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Readfile", err.Error())
		return nil, err
	}
	if err := json.Unmarshal(bytes, &xxx); err != nil {
		fmt.Println("Unmarshar", err.Error())
		return nil, err
	}
	return xxx, nil

}

func main() {
	xxxMap, err := readfile("./test1.js")
	if err != nil {
		fmt.Println(err.Error())

	}
	fmt.Println(xxxMap)

}
