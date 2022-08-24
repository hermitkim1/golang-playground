package aa

import "fmt"

func init() {
	fmt.Println("aa package's init() - Start ")

	fmt.Println("aa package's init() - End ")
}

func CallAA() {

	fmt.Println()
	fmt.Println("CallAA() - Start ")

	fmt.Println("I'm aa package")

	fmt.Println("CallAA() - End")
	fmt.Println()
}
