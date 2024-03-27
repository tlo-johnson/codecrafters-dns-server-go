package main

import (
	"encoding/binary"
)

type dnsAnswer struct {
  name []byte
  dnsType uint16
  class uint16
  timeToLive uint32
  length uint16
  data []byte
}

type dnsAnswers struct {
  answers []dnsAnswer
}

func (answer dnsAnswer) bytes() []byte {
  bytes := answer.name
  bytes = binary.BigEndian.AppendUint16(bytes, answer.dnsType)
  bytes = binary.BigEndian.AppendUint16(bytes, answer.class)
  bytes = binary.BigEndian.AppendUint32(bytes, answer.timeToLive)
  bytes = binary.BigEndian.AppendUint16(bytes, answer.length)
  bytes = append(bytes, answer.data...)

  return bytes
}

func newDnsAnswer(question dnsQuestion) dnsAnswer {
  questionBytes := question.bytes()
  questionLength := len(questionBytes)

  return dnsAnswer {
    name: questionBytes[: questionLength - 4],
    dnsType: 1,
    class: 1,
    timeToLive: 60,
    length: 4,
    data: []byte{ 8, 8, 8, 8 },
  }
}

func newDnsAnswers(questions dnsQuestions) dnsAnswers {
  var dnsAnswers dnsAnswers

  for _, question := range questions.questions {
    dnsAnswers.answers = append(dnsAnswers.answers, newDnsAnswer(question))
  }

  return dnsAnswers
}

func (answers dnsAnswers) bytes() []byte {
  var bytes []byte
  for _, answer := range answers.answers {
    bytes = append(bytes, answer.bytes()...)
  }
  return bytes
}
