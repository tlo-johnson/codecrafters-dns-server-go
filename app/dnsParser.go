package main

import (
	"encoding/binary"
)


func parseRequest(request []byte) (dnsHeader, dnsQuestions) {
  header := parseHeader(request)
  questions := parseQuestions(request, header.questionCount)

  return header, questions
}

func parseDnsFlags(data []byte) []byte {
  firstByte := 0b10000000 | data[2] & 0b11111001
  secondByte := 0b00001111 | data[3]
   if secondByte == 0 {
    return []byte { firstByte, 0 }
  } else {
    return []byte { firstByte, 4 }
  }
}


func parseHeader(data []byte) dnsHeader {
  flags := parseDnsFlags(data)

  return dnsHeader {
    packetIdentifier: binary.BigEndian.Uint16(data[0:2]),
    queryResponseIndicator: flags[0] >> 7,
    operationCode: (flags[0] & 0b01111000) >> 3,
    authoritativeAnswer: 0,
    truncation: 0,
    recursionDesired: flags[0] & 0b00000001,
    recursionAvailable: 0,
    reserved: 0,
    responseCode: flags[1],
    questionCount: binary.BigEndian.Uint16(data[4:6]),
    answerRecordCount: binary.BigEndian.Uint16(data[4:6]),
    authorityRecordCount: 0,
    additionalRecordCount: 0,
  }
}

func parseQuestions(data []byte, questionCount uint16) dnsQuestions {
  var dnsQuestions dnsQuestions

  index := 12
  var name []byte
  for questionCount > 0 {
    switch {
      case data[index] == 0:
        name = append(name, 0)
        dnsQuestion := newDnsQuestion(name)
        dnsQuestions.questions = append(dnsQuestions.questions, dnsQuestion)

        index += 4
        questionCount--
        name = make([]byte, 0)

      case data[index] & 0b11000000 != 0:
        offset := binary.BigEndian.Uint16(data[index : index + 2]) & 0b0011111111111111
        index = int(offset) - 1 // remember that index gets incremented once we get out the switch statement

      default:
        length := int(data[index])
        name = append(name, data[index : index + length + 1]...)
        index += length
    }

    index++
  }

  return dnsQuestions
}
