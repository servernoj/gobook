package main

import (
	"fmt"
	"os"
)

func main() {

	f := &flags{
		c: 10,
		n: 100,
	}
	if err := f.parse(); err != nil {
		os.Exit(1)
	}

	fmt.Printf("%+v\n", f)

}
