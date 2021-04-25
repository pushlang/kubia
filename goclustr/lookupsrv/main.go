package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	_, addrs, err := net.LookupSRV("", "", "kubia")
	if err != nil {
		log.Fatal(err)
	}

	for _, addr := range addrs {
		fmt.Printf("Target:%s Port: %d \n", addr.Target, addr.Port)
		port, err := net.LookupPort("tcp", "8080")

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Port:%d\n", port)
	}
}
