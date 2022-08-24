package bb

import "fmt"

func init() {
	fmt.Println("bb package's init() - Start ")

	fmt.Println("bb package's init() - End ")
}

func CallBB() {
	fmt.Println()
	fmt.Println("CallBB() - Start ")

	fmt.Println("I'm bb package")

	fmt.Println("CallBB() - End")
	fmt.Println()
}
