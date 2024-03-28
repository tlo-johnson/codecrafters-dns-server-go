package main

import (
	"encoding/binary"
)


func parseDnsMessage(request []byte) (dnsHeaderSection, dnsQuestionSection, dnsAnswerSection) {
  header := parseHeader(request)
  questions, answerSectionIndex := parseQuestions(request, header.questionCount)
  answers := parseAnswers(request, header.answerRecordCount, answerSectionIndex)

  return header, questions, answers
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


func parseHeader(data []byte) dnsHeaderSection {
  flags := parseDnsFlags(data)

  return dnsHeaderSection {
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
    answerRecordCount: binary.BigEndian.Uint16(data[6:8]),
    authorityRecordCount: 0,
    additionalRecordCount: 0,
  }
}

func parseQuestions(data []byte, questionCount uint16) (dnsQuestionSection, int) {
  var questions dnsQuestionSection

  index := 12
  var name []byte
  for questionCount > 0 {
    switch {
      case data[index] == 0:
        question, nextIndex := createQuestion(name, data, index)
        questions.questions = append(questions.questions, question)

        index = nextIndex
        questionCount--
        name = make([]byte, 0)

      case data[index] & 0b11000000 != 0:
        offset := binary.BigEndian.Uint16(data[index : index + 2]) & 0b0011111111111111
        index = int(offset)

      default:
        length := int(data[index])
        name = append(name, data[index : index + length + 1]...)
        index += length + 1
    }
  }

  return questions, index
}

func createQuestion(name []byte, data []byte, index int) (dnsQuestion, int) {
  name = append(name, data[index])
  index++

  dnsType := binary.BigEndian.Uint16(data[index:index + 2])
  index += 2

  class := binary.BigEndian.Uint16(data[index:index + 2])
  index += 2

  dnsQuestion := dnsQuestion {
    name: name,
    dnsType: dnsType,
    class: class,
  }

  return dnsQuestion, index
}

func parseAnswers(data []byte, answerCount uint16, index int) dnsAnswerSection {
  var answers dnsAnswerSection

  var name []byte
  for answerCount > 0 {
    switch {
      case data[index] == 0:
        answer, nextIndex := createAnswer(name, data, index)
        answers.answers = append(answers.answers, answer)

        index = nextIndex
        answerCount--
        name = make([]byte, 0)

      case data[index] & 0b11000000 != 0:
        offset := binary.BigEndian.Uint16(data[index : index + 2]) & 0b0011111111111111
        index = int(offset)

      default:
        length := int(data[index])
        name = append(name, data[index:index + length + 1]...)
        index += length + 1
    }
  }

  return answers
}

func createAnswer(name []byte, data []byte, index int) (dnsAnswer, int) {
  name = append(name, data[index])
  index++

  dnsType := binary.BigEndian.Uint16(data[index:index + 2])
  index += 2

  class := binary.BigEndian.Uint16(data[index:index + 2])
  index += 2

  timeToLive := binary.BigEndian.Uint32(data[index:index + 4])
  index += 4

  length := binary.BigEndian.Uint16(data[index:index + 2])
  index += 2

  rdata := data[index:index + int(length)]
  index += int(length)

  answer := dnsAnswer {
    name: name,
    dnsType: dnsType,
    class: class,
    timeToLive: timeToLive,
    length: length,
    data: rdata,
  }

  return answer, index
}
