package model

import "time"

type Message struct {
	Text        string
	CreatedAt   time.Time
	SourceBotId string
	TargetBotId string
}
