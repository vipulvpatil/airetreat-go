package model

import (
	"sort"
	"time"
)

type Message struct {
	Text        string
	CreatedAt   time.Time
	SourceBotId string
	TargetBotId string
	MessageType string
}

func (m *Message) IsQuestion() bool {
	return m.MessageType == "question"
}

func (m *Message) IsAnswer() bool {
	return m.MessageType == "answer"
}

type DetailedMessage struct {
	Text          string
	CreatedAt     time.Time
	SourceBotId   string
	SourceBotName string
	TargetBotId   string
	TargetBotName string
	MessageType   string
}

type detailedMessageSortByCreatedAt []DetailedMessage

func (m detailedMessageSortByCreatedAt) Len() int {
	return len(m)
}

func (m detailedMessageSortByCreatedAt) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m detailedMessageSortByCreatedAt) Less(i, j int) bool {
	return m[i].CreatedAt.Before(m[j].CreatedAt)
}

func (m detailedMessageSortByCreatedAt) sort() {
	sort.Sort(m)
}
