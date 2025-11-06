package main

import (
	"fmt"
	"net"
	"net/netip"
	"os"
	"strings"
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

func scanPorts(addr string, network string, duration time.Duration) {
	minPort := 1
	maxPort := 65535 //65535
	ports := make([]int, maxPort+1)
	//portsChan := make(chan []int)
	//defer close(portsChan)
	var wg sync.WaitGroup
	for i := minPort; i < maxPort; i++ {
		wg.Go(func() {
			//singlePortScan(fmt.Sprintf("127.0.0.1:%d", i), time.Second, i)
			if ok, _ := singlePortScan(addr, network, duration); ok {
				//ports = append(ports, i)
				//tempSlice := make([]int, 1)
				//tempSlice[0] = i
				//portsChan <- tempSlice
				ports[i] = 1
			} else {
				fmt.Println("Error")
			}
		})
		if i%10 == 0 {
			wg.Wait()
		}
	}
	wg.Wait()
	//fmt.Println(len(portsChan))
	//fmt.Println(ports[443])
	//return ports
	for i := 1; i < maxPort+1; i++ {
		if ports[i] == 1 {
			fmt.Println(i)
		}
	}
}

func scanAddresses(addrs *[]netip.Addr) {
	fmt.Println((*addrs)[0].String())
	/*
		var wg sync.WaitGroup
		for i := 0; i < len(*addrs); i++ {
			wg.Go(func() {
				if ok, _ := singlePortScan((*addrs)[i].String(), "tcp", time.Second); ok {
					fmt.Printf("Good: %s\n", (*addrs)[i].String())
				} else {
					fmt.Printf("Bad: %s\n", (*addrs)[i].String())
				}
			})

			if i%10 == 0 {
				wg.Wait()
			}
		}
		wg.Wait()
	*/
}

func portOpen(network string, addr string, duration time.Duration) bool {
	isOpen := true

	conn, err := net.DialTimeout(network, addr, duration)
	if conn != nil {
		conn.Close()
	}
	if err != nil {
		isOpen = false
	}

	return isOpen
}

func singlePortScan(addr string, network string, duration time.Duration) (bool, error) {
	success := true
	conn, err := net.DialTimeout(network, addr, duration)
	if conn != nil {
		conn.Close()
	}
	if err != nil {
		success = false
		if opError, ok := err.(*net.OpError); ok && !strings.Contains(opError.Err.Error(), "conn") {
			err = nil
		}
	}
	return success, err
}
