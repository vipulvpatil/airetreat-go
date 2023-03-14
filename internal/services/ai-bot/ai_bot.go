package aibot

import "github.com/vipulvpatil/airetreat-go/internal/model"

type AiQuestionGenerator interface {
	GetNextQuestion() string
}

type AiAnswerGenerator interface {
	GetNextAnswer() string
}

type aiBot struct {
	name              string
	conversationSoFar string
}

type AiBotOptions struct {
	Name string
	Bots []*model.Bot
}

func (ab *aiBot) GetNextQuestion() string {
	return randomFallbackQuestion()
}

func (ab *aiBot) GetNextAnswer() string {
	return randomFallbackAnswer()
}

func randomFallbackQuestion() string {
	return "What are we really talking about?"
}

func randomFallbackAnswer() string {
	return "I am unsure how to answer that"
}
