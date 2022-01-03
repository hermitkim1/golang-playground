package main

import "fmt"

func main() {

	placeHolderBody := `{"command": "%s", "userName": "cb-user"}`

	doubleBackslashed := fmt.Sprintf("[\"123.123.123.123:123\"]")
	fmt.Printf("body: %#v\n", doubleBackslashed)

	body := fmt.Sprintf(placeHolderBody, doubleBackslashed)
	fmt.Printf("body: %#v\n", body)

}
