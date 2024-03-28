package main

type dnsMessage struct {
  headerSection dnsHeaderSection
  questionSection dnsQuestionSection
  answerSection dnsAnswerSection
}

func (message dnsMessage) bytes() []byte {
  var result []byte
  result = append(result, message.headerSection.bytes()...)
  result = append(result, message.questionSection.bytes()...)
  result = append(result, message.answerSection.bytes()...)

  return result
}

func (request dnsMessage) retrieveAnswers(resolver string) dnsMessage {
  response := dnsMessage {
    headerSection: request.headerSection,
    questionSection: request.questionSection,
    answerSection: dnsAnswerSection { },
  }

  for _, question := range request.questionSection.questions {
    answer := forwardRequest(question, resolver)
    response.answerSection.answers = append(response.answerSection.answers, answer)
    response.headerSection.answerRecordCount++
  }

  return response
}

func forwardRequest(question dnsQuestion, resolver string) dnsAnswer {
  answers := dnsAnswerSection { }

  header := dnsHeaderSection { questionCount: 1 }
  questions := dnsQuestionSection {
    questions: []dnsQuestion { question },
  }
  message := dnsMessage { header, questions, answers }

  response := forwardMessage(message, resolver)
  return response.answers[0]
}
