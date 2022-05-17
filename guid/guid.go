package main

import (
	"fmt"

	"github.com/rs/xid"
)

func main() {
	guid := xid.New()

	fmt.Printf("guid.Machine(): %s\n", string(guid.Machine()))
	fmt.Printf("guid.String(): %s\n", guid.String())
	fmt.Printf("guid.Pid(): %d\n", guid.Pid())
	fmt.Printf("guid.Time(): %s\n", guid.Time())
	fmt.Printf("guid.Counter()%d\n", guid.Counter())
}
