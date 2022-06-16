package main

import (
	"fmt"

	"github.com/go-ping/ping"
)

func main() {
	desIP := "172.21.245.5"
	fmt.Printf("Ping to %s\n", desIP)

	pinger, err := ping.NewPinger(desIP)
	pinger.SetPrivileged(true)
	if err != nil {
		fmt.Println(err)
	}
	//pinger.OnRecv = func(pkt *ping.Packet) {
	//	fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
	//		pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	//}

	size := 256 // default 64 bytes
	pinger.Count = 10
	pinger.Size = size
	pinger.Run() // blocks until finished

	stats := pinger.Statistics() // get send/receive/rtt stats
	// outVal.MininumRTT = stats.MinRtt.Seconds()
	// outVal.AverageRTT = stats.AvgRtt.Seconds()
	// outVal.MaximumRTT = stats.MaxRtt.Seconds()
	// outVal.StdDevRTT = stats.StdDevRtt.Seconds()
	// outVal.PacketsReceive = stats.PacketsRecv
	// outVal.PacketsLoss = stats.PacketsSent - stats.PacketsRecv
	// outVal.BytesReceived = stats.PacketsRecv * size

	fmt.Printf("%v\n", stats)

	fmt.Printf("round-trip min/avg/max/stddev/dupl_recv = %v/%v/%v/%v/%v bytes\n",
		stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt, stats.PacketsRecv*size)
}
