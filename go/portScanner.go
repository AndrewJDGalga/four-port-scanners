package main

import (
	"fmt"
	"net"
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
	_, err := singlePortScan("127.0.0.1:80", "tcp", time.Second)
	if err != nil {
		fmt.Println(err)
	}
	_, err = singlePortScan("172.16.24.24:445", "tcp", time.Second)
	if err == nil {
		fmt.Println("no error")
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
