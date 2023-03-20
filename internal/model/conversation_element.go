package model

import (
	"sort"
	"time"
)

type ConversationElement struct {
	IsQuestion bool
	BotId      string
	BotName    string
	Text       string
}

type sortedConversationElement struct {
	IsQuestion bool
	BotId      string
	BotName    string
	Text       string
	CreatedAt  time.Time
}

type conversationSortByCreatedAt []sortedConversationElement

func (m conversationSortByCreatedAt) Len() int {
	return len(m)
}

func (m conversationSortByCreatedAt) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m conversationSortByCreatedAt) Less(i, j int) bool {
	return m[i].CreatedAt.Before(m[j].CreatedAt)
}

func (m conversationSortByCreatedAt) sortAndConvertToConversation() []ConversationElement {
	sort.Sort(m)
	conversation := []ConversationElement{}
	for _, elem := range m {
		conversation = append(conversation, ConversationElement{
			IsQuestion: elem.IsQuestion,
			BotId:      elem.BotId,
			BotName:    elem.BotName,
			Text:       elem.Text,
		})
	}
	return conversation
}
