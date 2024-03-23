package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

type dnsMessage struct {
  header []byte
  question []byte
  answer []byte
}

func newDnsMessage(data []byte) *dnsMessage {
  return &dnsMessage {
    header: newHeader(data),
    question: newQuestion(data),
    answer: newAnswer(data),
  }
}

func (message *dnsMessage) byte() []byte {
  result := append(message.header[:], message.question[:]...)
  result = append(result, message.answer[:]...)
  return result
}

func newHeader(data []byte) []byte {
  header := make([]byte, 12)
  copy(header, data[0:2])
  copy(header[2:4], newFlags(data))
  binary.BigEndian.PutUint16(header[4:6], 1) // Question Count (QDCOUNT)
  binary.BigEndian.PutUint16(header[6:8], 1) // Answer Record Count (ANCOUNT)

  return header
}

func newFlags(data []byte) []byte {
  firstByte := 0b10000000 | data[2] & 0b11111001
  secondByte := 0b00001111 | data[3]
   if secondByte == 0 {
    return []byte { firstByte, 0 }
  } else {
    return []byte { firstByte, 4 }
  }
}

func newQuestion(data []byte) []byte {
  var question []byte

  for _, value := range data[12:] {
    if value == 0 {
      break
    }
    question = append(question, value)
  }

  question = append(question, 0)
  question = binary.BigEndian.AppendUint16(question, 1) // "A" record
  question = binary.BigEndian.AppendUint16(question, 1) // "IN" record

  return question
}

func newAnswer(data []byte) []byte {
  answer := domainName(data)
  answer = binary.BigEndian.AppendUint16(answer, 1)
  answer = binary.BigEndian.AppendUint16(answer, 1)
  answer = binary.BigEndian.AppendUint32(answer, 1)
  answer = binary.BigEndian.AppendUint16(answer, 4)
  answer = append(answer, []byte{ 8, 8, 8, 8 }...)

  return answer
}

func domainName(data []byte) []byte {
  var result []byte

  for _, value := range data[12:] {
    if value == 0 {
      result = append(result, 0)
      break
    }
    result = append(result, value)
  }

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
	
		receivedData := buf[:size]
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, string(receivedData))
	
    response := newDnsMessage(receivedData)
	
    _, err = udpConn.WriteToUDP(response.byte(), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
