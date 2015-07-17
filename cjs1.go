package main

import (
	"fmt"
	"os/exec"
	//"reflect"
	"regexp"
)

func main() {
	var a1, a2 string
	argv1 := []string{"-n", "2p", "./server_config_CN_PROD.js"}
	a := exec.Command("sed", argv1...)
	argv2 := []string{"-n", "2,10000p", "./server_config_CN_PROD.js"}
	b := exec.Command("sed", argv2...)
	argv3 := []string{"-n", "1p", "./server_config_CN_PROD.js"}
	b1 := exec.Command("sed", argv3...)
	d1, _ := a.Output()
	fmt.Println(string(d1))
	c := string(d1)

	d2, _ := b.Output()
	fmt.Println(string(d2))
	d3, _ := b1.Output()

	filedata1, _ := regexp.Compile("nba_redis102")
	n1 := "nba_redis103"
	a1 = filedata1.ReplaceAllString(c, n1)
	filedata2, _ := regexp.Compile("_id:102")
	n2 := "_id:103"
	a2 = filedata2.ReplaceAllString(a1, n2)

	fmt.Println(a1)
	fmt.Println(a2)
	fmt.Println(string(d3))

	new := d3 + a2 + d2
	fmt.Println(new)
}
