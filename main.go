package main

import (
	"errors"
	"fmt"

	"github.com/Yves848/wingettui/libs"
	"github.com/Yves848/wingettui/winget"
)

func main() {
	SayHello()
	libs.SayHello2()
	out, err := winget.Invoke("list")
	if err != nil {
		fmt.Println(errors.New("error invoking winget"))
	}
	fmt.Println(string(out))

}
