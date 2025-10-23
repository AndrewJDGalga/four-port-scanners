package main

import (
	"fmt"
	"net"
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
	singlePortScanErr("172.16.24.24:445", time.Second, 445) //dial tcp 172.16.24.24:445: i/o timeout
	singlePortScanErr("127.0.0.1:80", time.Second, 80)
	//singlePortScan("127.0.0.1:80", time.Second, 80)      //dial tcp 127.0.0.1:80: connectex: No connection could be made because the target machine actively refused it.
}

func scanPorts() {
	var wg sync.WaitGroup
	for i := 1; i < 65535; i++ {
		wg.Go(func() {
			singlePortScan(fmt.Sprintf("127.0.0.1:%d", i), time.Second, i)
		})
		if i%10 == 0 {
			wg.Wait()
		}
	}
	wg.Wait()
}

func singlePortScanErr(addrAndPort string, duration time.Duration, port int) {
	_, err := net.DialTimeout("tcp", addrAndPort, duration)
	if err != nil {
		if opError, ok := err.(*net.OpError); ok {

			switch opError.Err.Error() {
			case "i/o timeout":
				fmt.Println("i o error")
			case "connectex: No connection could be made because the target machine actively refused it.":
				fmt.Println("connection error")
			default:
				fmt.Println("no conditions triggered")
			}

			//fmt.Println(opError.Err.Error())
		}

	}
}

func singlePortScan(addrAndPort string, duration time.Duration, port int) {
	conn, err := net.DialTimeout("tcp", addrAndPort, duration)
	if err == nil {
		conn.Close()
		fmt.Printf("%d Open\n", port)
	} else {
		fmt.Printf("Error: '%v'\nOf type '%T'\n", err, err)

		if conn != nil {
			conn.Close()
		}
	}
}
