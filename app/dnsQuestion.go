package main

import "encoding/binary"

type dnsQuestion struct {
  name []byte
  dnsType uint16
  class uint16
}

func (question dnsQuestion) bytes() []byte {
  bytes := question.name
  bytes = binary.BigEndian.AppendUint16(bytes, question.dnsType)
  bytes = binary.BigEndian.AppendUint16(bytes, question.class)

  return bytes
}

type dnsQuestionSection struct {
  questions []dnsQuestion
}

func (questions dnsQuestionSection) bytes() []byte {
  var bytes []byte
  for _, question := range questions.questions {
    bytes = append(bytes, question.bytes()...)
  }
  return bytes
}
