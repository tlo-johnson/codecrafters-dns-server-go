package main

import (
	"fmt"
	"net"
)

func forwardMessage(message dnsMessage, resolver string) dnsAnswerSection {
  conn, err := net.Dial("udp", resolver)
  if err != nil {
    fmt.Println("an error occurred connecting to forwarding server", err)
    return dnsAnswerSection { }
  }

  defer conn.Close()

  conn.Write(message.bytes())

  response := make([]byte, 2048)
  size, err := conn.Read(response)
  if err != nil {
    fmt.Println("an error occurred reading response from forwarding server", err)
    return dnsAnswerSection { }
  }

  _, _, answers := parseDnsMessage(response[:size])
  return answers
}
