/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/lqwangxg/myapp/cmd"
)

func main() {
	cmd.Execute()
	// for {
	// 	reader := bufio.NewReader(os.Stdin)
	// 	fmt.Print("> ")
	// 	_, err := reader.ReadString('\n')
	// 	if err != nil {
	// 		return
	// 	}
	// 	cmd.Execute()
	// }
}
