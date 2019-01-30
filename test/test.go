package main

import (
	"fmt"
	"github.com/hlfstr/aeridya"
	"os"
)

func main() {
	a, _, e := aeridya.Create("./conf")
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
	e = a.Run()
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}
