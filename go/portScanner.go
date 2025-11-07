package main

import (
	"fmt"
	"net"
	"net/netip"
	"os"
	"sync"
	"time"
)

func main() {
	if len(os.Args) > 1 {
		args := os.Args[1:]
		tempAddr, err := netip.ParseAddr(args[0])
		if err != nil {
			fmt.Printf("Address must be in acceptable IPv4/IPv6 format without port numbers. Examples:\n0.0.0.0\n::1")
			return
		}
		//addr := []netip.Addr{tempAddr}
		//scanAddresses(&addr)
		scanAddresses(&[]netip.Addr{tempAddr})
	} else {
		fmt.Println("Correct usage:\n\t>go run portScanner.go [ip_address(:optional_port)(, additional_ip_addresses)]")
	}
}

func scanPorts(addr string, network string, duration time.Duration) *[]int {
	minPort := 1
	maxPort := 65535
	ports := make([]int, maxPort+1)
	var wg sync.WaitGroup

	//portChan := make(chan []int, maxPort+1)

	for i := minPort; i < maxPort; i++ {
		wg.Go(func() {
			if portOpen(addr+fmt.Sprintf(":%d", i), network, duration) {
				ports[i] = i
				fmt.Println("Val: ", i)
			}
		})
		if i%10 == 0 {
			wg.Wait()
		}
	}
	wg.Wait()
	for i := 1; i < maxPort+1; i++ {
		if ports[i] == 1 {
			fmt.Println(i)
		}
	}
	return &ports
}

func scanAddresses(addrs *[]netip.Addr) {
	var wg sync.WaitGroup
	for i := 0; i < len(*addrs); i++ {
		wg.Go(func() {
			udpPorts := scanPorts((*addrs)[i].String(), "udp", time.Second)
			//tcpPorts := scanPorts((*addrs)[i].String(), "tcp", time.Second)

			fmt.Println((*udpPorts)[0])
			/*
				for i := 0; i < len(*udpPorts); i++ {
					if (*udpPorts)[i] != 0 {
						fmt.Println((*udpPorts)[i])
					}
				}*/

			/*
				if singlePortScan((*addrs)[i].String()+":443", "tcp", time.Second) {
					fmt.Printf("Good: %s\n", (*addrs)[i].String())
				} else {
					fmt.Printf("Bad: %s\n", (*addrs)[i].String())
				}
				if singlePortScan((*addrs)[i].String()+":443", "udp", time.Second) {
					fmt.Printf("Good: %s\n", (*addrs)[i].String())
				} else {
					fmt.Printf("Bad: %s\n", (*addrs)[i].String())
				}
			*/
		})

		if i%10 == 0 {
			wg.Wait()
		}
	}
	wg.Wait()
}

func portOpen(addr string, network string, duration time.Duration) bool {
	success := true
	conn, err := net.DialTimeout(network, addr, duration)
	if conn != nil {
		conn.Close()
	}
	if err != nil {
		success = false
		err = nil //good not to leave floating around?
	}
	return success
}
