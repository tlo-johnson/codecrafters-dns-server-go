package main

import "encoding/binary"

type dnsQuestion struct {
  name []byte
  dnsType uint16
  class uint16
}

func newDnsQuestion(name []byte) dnsQuestion {
  return dnsQuestion {
    name: name,
    dnsType: 1,
    class: 1,
  }
}

func (question dnsQuestion) bytes() []byte {
  bytes := question.name
  bytes = binary.BigEndian.AppendUint16(bytes, question.dnsType)
  bytes = binary.BigEndian.AppendUint16(bytes, question.class)

  return bytes
}

type dnsQuestions struct {
  questions []dnsQuestion
}

func (questions dnsQuestions) bytes() []byte {
  var bytes []byte
  for _, question := range questions.questions {
    bytes = append(bytes, question.bytes()...)
  }
  return bytes
}
