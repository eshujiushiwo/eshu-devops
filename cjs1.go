package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	//"reflect"
)

func printlog(c2 *exec.Cmd) {

	stdout, err := c2.StdoutPipe()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err := c2.Start(); err != nil {
		fmt.Println("Command  error:", err.Error())
		os.Exit(1)
	}
	in := bufio.NewScanner(stdout)
	fmt.Println(in)

	for in.Scan() {

		fmt.Println(in.Text())

	}
	if err := in.Err(); err != nil {
		fmt.Println("Err:", err.Error())
		os.Exit(1)
	}
}

func main() {

	argv1 := []string{"./LICENSE"}
	c1 := exec.Command("cat", argv1...)

	printlog(c1)
	/*stdout, err := c1.StdoutPipe()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err := c1.Start(); err != nil {
		fmt.Println("Command  error:", err.Error())
		os.Exit(1)
	}
	in := bufio.NewScanner(stdout)

	for in.Scan() {
		fmt.Println("9")
		fmt.Println(in.Text())
		fmt.Println("4")
	}
	if err := in.Err(); err != nil {
		fmt.Println("Err:", err.Error())
		os.Exit(1)
	}
	*/

}
