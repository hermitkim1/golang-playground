package main

import (
	"fmt"

	"github.com/yunkon-kim/golang-playground/init-sequence/pkg/aa"
	"github.com/yunkon-kim/golang-playground/init-sequence/pkg/bb"
	"github.com/yunkon-kim/golang-playground/init-sequence/pkg/cc"
)

func init() {
	fmt.Println("main package's init() - Start ")

	fmt.Println("main package's init() - End ")
}

func main() {

	fmt.Println()
	fmt.Println("main() - Start ")

	cc.CallCC()
	bb.CallBB()
	aa.CallAA()

	fmt.Println("main() - End")
	fmt.Println()
}
