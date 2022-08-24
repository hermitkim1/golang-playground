package cc

import (
	"fmt"

	"github.com/hermitkim1/golang-playground/init-sequence/pkg/bb"
)

func init() {
	fmt.Println("cc package's init() - Start ")

	fmt.Println("cc package's init() - End ")
}

func CallCC() {

	fmt.Println()
	fmt.Println("Call() - Start ")

	fmt.Println("I'm cc package")
	bb.CallBB()

	fmt.Println("CallCC() - End")
	fmt.Println()
}
