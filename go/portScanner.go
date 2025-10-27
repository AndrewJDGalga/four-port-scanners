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

	//scanAddresses()
	scanSubnetNetIP()
}

func scanSubnetNetIP() {
	testNet := "10.10.0.0"
	//testNet := "64:ff9b:1::/48"
	//testNet := "fe80::1ff:fe23:4567:890a%3"

	//prefix, _ := netip.ParsePrefix(testNet)
	address, _ := netip.ParseAddr(testNet)
	prefix := netip.PrefixFrom(address, 16)
	//fmt.Println(prefix.Addr().Next())
	//fmt.Println(prefix.Contains(prefix.Addr().Next()))
	//fmt.Println(prefix.Addr())

	//oprefix := netip.PrefixFrom("10.0.0.0", 16)

	for outOfRange := true; outOfRange; outOfRange = prefix.Contains(prefix.Addr().Next()) {
		fmt.Println(prefix.Addr())
		prefix = netip.PrefixFrom(prefix.Addr().Next(), 16)
	}
}

func scanPorts() {
	var wg sync.WaitGroup
	for i := 1; i < 65535; i++ {
		wg.Go(func() {
			//singlePortScan(fmt.Sprintf("127.0.0.1:%d", i), time.Second, i)
		})
		if i%10 == 0 {
			wg.Wait()
		}
	}
	wg.Wait()
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
