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
	var wg sync.WaitGroup
	for i := 1; i < 1000; i++ {
		wg.Go(func() {
			singlePortScan(fmt.Sprintf("127.0.0.1:%d", i), time.Second, i)
		})
		if i%10 == 0 {
			wg.Wait()
		}
	}

}

func singlePortScan(addrAndPort string, duration time.Duration, port int) {
	conn, err := net.DialTimeout("tcp", addrAndPort, duration)
	if err == nil {
		conn.Close()
		fmt.Printf("%d Open\n", port)
	} else {
		if conn != nil {
			fmt.Println(err)
			conn.Close()
		}
	}
}
