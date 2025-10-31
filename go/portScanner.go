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
	/*
		addPtr := flag.String("a", "127.0.0.1", "Address(es) for scanning. Defaults to loopback.")
		portRngPtr := flag.String("p", "1,65535", "Port(s) range for scanning. Defaults to 1,65535 == starting at 1, ending at 65535")
		if len(os.Args) > 1 {
			fmt.Println(os.Args)
		}
	*/
	//singlePortScanErr("172.16.24.24:445", time.Second, 445) //dial tcp 172.16.24.24:445: i/o timeout
	//singlePortScanErr("127.0.0.1:80", time.Second, 80)
	//singlePortScan("127.0.0.1:80", time.Second, 80)      //dial tcp 127.0.0.1:80: connectex: No connection could be made because the target machine actively refused it.
	/*
		_, err := singlePortScan("127.0.0.1:80", "tcp", time.Second)
		if err != nil {
			fmt.Println(err)
		}
		_, err = singlePortScan("172.16.24.24:445", "tcp", time.Second)
		if err == nil {
			fmt.Println("no error")
		}
	*/

	//testNet := "64:ff9b:1::/48"
	//testNet := "fe80::1ff:fe23:4567:890a%3"
	//scanSubnet("192.168.1.10", 26)
	//scanPorts("192.168.1.66", "tcp", time.Second)
	//fmt.Println(len(slicePtr))
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
