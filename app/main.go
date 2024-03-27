package main

import (
	"flag"
	"fmt"
	"net"
)

func respond(request []byte) dnsMessage {
  header, questions := parseRequest(request)
  return newDnsMessage(header, questions)
}

var resolver string

func parseFlags() {
  flag.StringVar(&resolver, "resolver", "", "<ip>:<port> for resolver")
  flag.Parse()
}

func main() {
  parseFlags()

	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}
	
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()
	
	buf := make([]byte, 512)
	
	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}
	
		request := buf[:size]
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, string(request))
	
    response := respond(request)
	
    _, err = udpConn.WriteToUDP(response.byte(), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
