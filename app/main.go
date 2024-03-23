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
  answer []byte
}

func newDnsMessage() *dnsMessage {
  return &dnsMessage {
    header: newHeader(),
    question: newQuestion(),
    answer: newAnswer(),
  }
}

func (message *dnsMessage) byte() []byte {
  result := append(message.header[:], message.question[:]...)
  result = append(result, message.answer[:]...)
  return result
}

func newHeader() []byte {
  header := make([]byte, 12)
  binary.BigEndian.PutUint16(header[0:2], 1234)
  binary.BigEndian.PutUint16(header[2:4], 1 << 15)
  binary.BigEndian.PutUint16(header[4:6], 1) // Question Count (QDCOUNT)
  binary.BigEndian.PutUint16(header[6:8], 1) // Answer Record Count (ANCOUNT)

  return header
}

func newQuestion() []byte {
  domainName := "codecrafters.io"

  question := labelSequence(domainName)
  question = binary.BigEndian.AppendUint16(question, 1) // "A" record
  question = binary.BigEndian.AppendUint16(question, 1) // "IN" record

  return question
}

func newAnswer() []byte {
  domainName := "codecrafters.io"

  answer := labelSequence(domainName)
  answer = binary.BigEndian.AppendUint16(answer, 1)
  answer = binary.BigEndian.AppendUint16(answer, 1)
  answer = binary.BigEndian.AppendUint32(answer, 1)
  answer = binary.BigEndian.AppendUint16(answer, 4)
  answer = append(answer, []byte{ 8, 8, 8, 8 }...)

  return answer
}

func labelSequence(domainName string) []byte {
  var result []byte
  parts := strings.Split(domainName, ".")
  for _, content := range(parts) {
    result = append(result, byte(len(content)))
    result = append(result, content...)
  }
  result = append(result, "\x00"...)

  return result
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
