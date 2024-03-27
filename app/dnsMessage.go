package main

type dnsMessage struct {
  header dnsHeader
  questions dnsQuestions
  answers dnsAnswers
}

func (message dnsMessage) byte() []byte {
  var result []byte
  result = append(result, message.header.byte()...)
  result = append(result, message.questions.bytes()...)
  result = append(result, message.answers.bytes()...)

  return result
}

func newDnsMessage(header dnsHeader, questions dnsQuestions) dnsMessage {
  return dnsMessage {
    header: header,
    questions: questions,
    answers: newDnsAnswers(questions),
  }
}
