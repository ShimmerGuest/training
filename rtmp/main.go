package main

import (
	`fmt`
)

func main() {
	err := Run(":1935")
	if err != nil {
		fmt.Println("err = ", err.Error())
	}
}
