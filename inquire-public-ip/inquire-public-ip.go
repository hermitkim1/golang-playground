package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

func main() {

	// [Warning] Occasionally fail to acquire public IP "https://ifconfig.co/"
	// The links below have not been tested.
	urls := []string{"https://ifconfig.co/",
		"https://api.ipify.org?format=text",
		"https://www.ipify.org",
		"http://myexternalip.com",
		"http://api.ident.me",
		"http://whatismyipaddress.com/api",
	}

	myPublicIP := ""

	for _, url := range urls {

		// Inquire public IP address
		fmt.Printf("URL: %s\n", url)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Print(err)
		}

		// Perform error handling
		defer func() {
			errClose := resp.Body.Close()
			if errClose != nil {
				fmt.Print("can't close the response", errClose)
			}
		}()

		// 결과 출력
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Print(err)
		}

		trimmed := strings.TrimSuffix(string(data), "\n") // Remove '\n' if exist
		fmt.Printf("Returned: %s", trimmed)

		if net.ParseIP(trimmed) != nil {
			myPublicIP = trimmed
			break
		}
	}

	fmt.Printf("Public IP address: %s", myPublicIP)

}
