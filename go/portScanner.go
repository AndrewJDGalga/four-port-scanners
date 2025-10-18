package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	addPtr := flag.String("a", "127.0.0.1", "Address(es) for scanning. Defaults to loopback.")
	portRngPtr := flag.String("p", "1,65535", "Port(s) range for scanning. Defaults to 1,65535 == starting at 1, ending at 65535")
	if len(os.Args) > 1 {
		fmt.Println(os.Args)
	}
}
