package main

import (
	"encoding/binary"
  "strings"
	"fmt"
	"net"
)

type dnsMessage struct {
  header []byte
  question []byte
}

func newDnsMessage() *dnsMessage {
  return &dnsMessage {
    header: newHeader(),
    question: newQuestion(),
  }
}

func (message *dnsMessage) byte() []byte {
  return append(message.header[:], message.question[:]...)
}

func newHeader() []byte {
  header := make([]byte, 12)
  binary.BigEndian.PutUint16(header[0:2], 1234)
  binary.BigEndian.PutUint16(header[2:4], 1 << 15)
  binary.BigEndian.PutUint16(header[4:6], 1)

  return header
}

func newQuestion() []byte {
  var question []byte

  domainName := "codecrafters.io"
  parts := strings.Split(domainName, ".")
  for _, content := range(parts) {
    question = append(question, byte(len(content)))
    question = append(question, content...)
  }
  question = append(question, "\x00"...)
  question = binary.BigEndian.AppendUint16(question, 1) // "A" record
  question = binary.BigEndian.AppendUint16(question, 1) // "IN" record

  return question
}

func main() {
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
	
		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)
	
    response := newDnsMessage()
	
    _, err = udpConn.WriteToUDP(response.byte(), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
