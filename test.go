package main

import (
	"fmt"
	"os"
	"os/exec"
)

func chdir(dir string) {
	err := os.Chdir(dir)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(os.Getwd())
}
func git_pull() {
	argv := []string{"pull"}
	c := exec.Command("git", argv...)
	d, _ := c.Output()
	fmt.Println(string(d))
}

func main() {

	chdir("/Users/zhou.liyang/eshu-devops/eshu-devops")
	git_pull()

	/*
		fmt.Println(os.Chdir("/Users/zhou.liyang"))
		fmt.Println(os.Getwd())

		argv := []string{"-l"}
		c := exec.Command("ls", argv...)
		d, _ := c.Output()
		fmt.Println(string(d))
	*/
}
