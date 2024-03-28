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

type dnsAnswerSection struct {
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

func (answers dnsAnswerSection) bytes() []byte {
  var bytes []byte
  for _, answer := range answers.answers {
    bytes = append(bytes, answer.bytes()...)
  }
  return bytes
}
