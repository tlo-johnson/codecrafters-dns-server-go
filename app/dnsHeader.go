package main

import "encoding/binary"

type dnsHeader struct {
  packetIdentifier uint16
  queryResponseIndicator uint8
  operationCode uint8
  authoritativeAnswer uint8
  truncation uint8
  recursionDesired uint8
  recursionAvailable uint8
  reserved uint8
  responseCode uint8
  questionCount uint16
  answerRecordCount uint16
  authorityRecordCount uint16
  additionalRecordCount uint16
}

func (header *dnsHeader) byte() []byte {
  var bytes []byte

  flags := uint16(header.queryResponseIndicator) << 15 | 
    uint16(header.operationCode) << 11 |
    uint16(header.authoritativeAnswer) << 10 |
    uint16(header.truncation) << 9 |
    uint16(header.recursionDesired) << 8 |
    uint16(header.recursionAvailable) << 7 |
    uint16(header.reserved) << 6 |
    uint16(header.responseCode)

  bytes = binary.BigEndian.AppendUint16(bytes, header.packetIdentifier)
  bytes = binary.BigEndian.AppendUint16(bytes, flags)
  bytes = binary.BigEndian.AppendUint16(bytes, header.questionCount)
  bytes = binary.BigEndian.AppendUint16(bytes, header.answerRecordCount)
  bytes = binary.BigEndian.AppendUint16(bytes, header.authorityRecordCount)
  bytes = binary.BigEndian.AppendUint16(bytes, header.additionalRecordCount)

  return bytes
}
