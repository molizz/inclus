package main

import (
	"fmt"

	"github.com/molizz/inclus/bin/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}
