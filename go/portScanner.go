package main

import (
	"fmt"
	"net"
	"net/netip"
	"strings"
	"sync"
	"time"
)

func main() {
	//fmt.Println(addressValid("192.168.1.66:443", 2*time.Second))
	pingAddress()
}

func scanSubnet(address string, subnet int) {
	addr, _ := netip.ParseAddr(address)
	prefix := netip.PrefixFrom(addr, subnet)

	for outOfRange := true; outOfRange; outOfRange = prefix.Contains(prefix.Addr().Next()) {
		//fmt.Println(prefix.Addr())
		//if singlePortScan(prefix.Addr().String(), prefix.String(), time.Second)

		prefix = netip.PrefixFrom(prefix.Addr().Next(), subnet)
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

func scanAddresses() {
	testAddresses := [5]string{"127.0.0.1:443", "172.16.24.24:443", "172.16.0.1:443", "1.1.1.1:443", "8.8.8.8:443"}

	var wg sync.WaitGroup
	for i := 0; i < len(testAddresses); i++ {
		wg.Go(func() {
			if ok, _ := singlePortScan(testAddresses[i], "tcp", time.Second); ok {
				fmt.Printf("Good: %s\n", testAddresses[i])
			} else {
				fmt.Printf("Bad: %s\n", testAddresses[i])
			}
		})

		if i%10 == 0 {
			wg.Wait()
		}
	}
	wg.Wait()
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

func pingAddress() {
	conn, err := net.DialTimeout("ip:icmp", "127.0.0.1", 2*time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}
	if conn != nil {
		fmt.Println("Pinged?")
		conn.Close()
	}
}

/*
func addressValid(addr string, duration time.Duration) bool {
	isValid := false

	conn, err := net.DialTimeout("tcp", addr, duration)
	if conn != nil {
		isValid = true
		conn.Close()
	}
	if err != nil {
		fmt.Println(err)
		if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "connection refused" {
			isValid = true
		}
	}

	return isValid
}
*/

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
