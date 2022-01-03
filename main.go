package main

import (
	"fmt"
	"time"

	"github.com/rs/xid"
)

func main() {
	guid := xid.New()

	println(guid.Counter())
	a := 0
	b := 0
	c := 0
	d := 0
	e := 0

	for i := 0; i < 1000000; i++ {
		// fmt.Println(t)
		now := time.Now()
		nano := now.Nanosecond()
		millisec := nano / 1000000
		// fmt.Println(t)
		switch millisec % 5 {
		case 0:
			a += 1
		case 1:
			b += 1
		case 2:
			c += 1
		case 3:
			d += 1
		case 4:
			e += 1
		}
		// time.Sleep(100 * time.Nanosecond)
	}

	sum := float64(a + b + c + d + e)
	fmt.Printf("%d / %d / %d / %d / %d\n", a, b, c, d, e)
	fmt.Printf("%.2f / %.2f / %.2f / %.2f / %.2f\n", float64(a)/sum, float64(b)/sum, float64(c)/sum, float64(d)/sum, float64(e)/sum)
}
