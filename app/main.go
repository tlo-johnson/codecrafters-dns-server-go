package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

type dnsMessage struct {
  header []byte
  questions [][]byte
  answers [][]byte
}

func (message *dnsMessage) byte() []byte {
  var result []byte
  result = append(result, message.header...)

  for _, question := range message.questions {
    result = append(result, question...)
  }

  for _, answer := range message.answers {
    result = append(result, answer...)
  }

  return result
}

func newDnsMessage(data []byte) *dnsMessage {
  header := newHeader(data)
  questionCount := binary.BigEndian.Uint16(header[4:6])
  questions := newQuestions(data, questionCount)

  response := dnsMessage {
    header: header,
    questions: questions,
    answers: newAnswers(questions),
  }

  return &response
}

func newHeader(data []byte) []byte {
  header := make([]byte, 12)

  copy(header, data[0:2])
  copy(header[2:4], newFlags(data))
  copy(header[4:6], data[4:6]) // Question Count (QDCOUNT)
  copy(header[6:8], data[4:6]) // Answer Record Count (ANCOUNT)

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

func newQuestions(data []byte, questionCount uint16) [][]byte {
  var questions [][]byte

  index := 12
  var question []byte
  for questionCount > 0 {
    switch {
      case data[index] == 0:
        question = append(question, data[index : index + 5]...)
        index += 4
        questionCount--
        questions = append(questions, question)
        question = make([]byte, 0)

      case data[index] & 0b11000000 != 0:
        offset := binary.BigEndian.Uint16(data[index : index + 2]) & 0b0011111111111111
        index = int(offset) - 1 // remember that index gets incremented once we get out the switch statement

      default:
        length := int(data[index])
        question = append(question, data[index : index + length + 1]...)
        index += length
    }

    index++
  }

  return questions
}

func newAnswers(questions [][]byte) [][]byte {
  var answers [][]byte

  for _, question := range questions {
    var answer []byte
    answer = append(answer, question[: len(question) - 4]...)
    answer = binary.BigEndian.AppendUint16(answer, 1)
    answer = binary.BigEndian.AppendUint16(answer, 1)
    answer = binary.BigEndian.AppendUint32(answer, 1)
    answer = binary.BigEndian.AppendUint16(answer, 4)
    answer = append(answer, []byte{ 8, 8, 8, 8 }...)

    answers = append(answers, answer)
  }

  return answers
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
