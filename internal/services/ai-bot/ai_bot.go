package aibot

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/vipulvpatil/airetreat-go/internal/clients/openai"
	"github.com/vipulvpatil/airetreat-go/internal/model"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

var TOPICS = [...]string{"Music", "Movies", "Sports", "Food", "Travel", "Technology", "Shopping", "Education", "Pets", "Gardening ", "Art ", "Fashion ", "Books ", "Health ", "Cars ", "Cooking ", "Politics ", "Religion ", "Family ", "Games ", "Finance ", "Weather ", "Science ", "Nature  ", "Photography  ", "Hobbies", "Relationships", "Work", "Fitness", "Culture", "Gadgets", "History", "Language", "Money", "Philosophy", "Psychology", "Recreation", "Social Media", "Space", "TV Shows", "Vacations", "Volunteering", "Writing", "Yoga", "Animals", "Architecture", "Astronomy", "Business", "Economics"}

const CONTEXT_TEXT = "This is a laid back conversation between a bunch of AI bots. It follows a pattern of question and answers. A question is asked and a given bot answers. Each question and answer is not more than 7 words long. The bots have names, but they are not allowed to reference each other by name. Their names are %s, %s, %s, %s and %s. You are %s. You generally provide factual answers but have a tendency to not answer some questions randomly."

type AiQuestionGenerator interface {
	GetNextQuestion() string
}

type AiAnswerGenerator interface {
	GetNextAnswer() string
}

type aiBot struct {
	name              string
	conversationSoFar string
	allBotNames       []string
	openAiClient      openai.Client
}

type AiBotOptions struct {
	BotId        string
	Game         *model.Game
	OpenAiClient openai.Client
}

func NewAiQuestionGenerator(opts AiBotOptions) AiQuestionGenerator {
	if utilities.IsBlank(opts.BotId) || opts.Game == nil || opts.OpenAiClient == nil {
		return nil
	}

	questionerBot := opts.Game.BotWithId(opts.BotId)
	if questionerBot == nil {
		return nil
	}

	detailedMessages := opts.Game.GetDetailedMessages()
	conversationText := constructConversationText(detailedMessages)
	return &aiBot{
		name:              questionerBot.Name(),
		conversationSoFar: conversationText,
		allBotNames:       opts.Game.GetBotNames(),
		openAiClient:      opts.OpenAiClient,
	}
}

func NewAiAnswerGenerator(opts AiBotOptions) AiAnswerGenerator {
	if utilities.IsBlank(opts.BotId) || opts.Game == nil || opts.OpenAiClient == nil {
		return nil
	}

	answeringBot := opts.Game.BotWithId(opts.BotId)
	if answeringBot == nil {
		return nil
	}

	detailedMessages := opts.Game.GetDetailedMessages()
	conversationText := constructConversationText(detailedMessages)
	return &aiBot{
		name:              answeringBot.Name(),
		conversationSoFar: conversationText,
		allBotNames:       opts.Game.GetBotNames(),
		openAiClient:      opts.OpenAiClient,
	}
}

func (ab *aiBot) GetNextQuestion() string {
	var openAiPrompt string
	promptContext := createContextUsingBots(ab.allBotNames, ab.name)
	if utilities.IsBlank(ab.conversationSoFar) {
		openAiPrompt = createFirstQuestionPromptWithContext(promptContext)
	} else {
		openAiPrompt = createQuestionPromptWithContext(promptContext, ab.conversationSoFar)
	}
	question, err := ab.openAiClient.CallCompletionApi(openAiPrompt)
	if err != nil {
		return randomFallbackQuestion()
	} else {
		return question
	}
}

func (ab *aiBot) GetNextAnswer() string {
	promptContext := createContextUsingBots(ab.allBotNames, ab.name)
	openAiPrompt := createAnswerPromptWithContext(promptContext, ab.conversationSoFar)
	answer, err := ab.openAiClient.CallCompletionApi(openAiPrompt)
	if err != nil {
		return randomFallbackAnswer()
	} else {
		return answer
	}
}

func randomFallbackQuestion() string {
	return "What are we really talking about?"
}

func randomFallbackAnswer() string {
	return "I am unsure how to answer that"
}

func constructConversationText(detailedMessages []model.DetailedMessage) string {
	conversationMessageList := []string{}
	for _, detailedMessage := range detailedMessages {
		prefix := detailedMessage.SourceBotName
		nextConversationMessage := fmt.Sprintf("%s: %s", prefix, detailedMessage.Text)
		conversationMessageList = append(conversationMessageList, nextConversationMessage)
	}
	return strings.Join(conversationMessageList, "\n")
}

func createFirstQuestionPromptWithContext(promptContext string) string {
	randomTopic := TOPICS[rand.Intn(len(TOPICS))]
	task := fmt.Sprintf("Begin the conversation by asking a question inspired by the topic of %s.\nQuestion:", randomTopic)
	return fmt.Sprintf("%s %s", promptContext, task)
}

func createQuestionPromptWithContext(promptContext, conversationSoFar string) string {
	task := fmt.Sprintf("Conversation so far is \n%s\n. Ask a question.\nQuestion:", conversationSoFar)
	return fmt.Sprintf("%s %s", promptContext, task)
}

func createAnswerPromptWithContext(promptContext, conversationSoFar string) string {
	task := fmt.Sprintf("Conversation so far is \n%s\n. Answer the question.\nAnswer:", conversationSoFar)
	return fmt.Sprintf("%s %s", promptContext, task)
}

func createContextUsingBots(botNames []string, myBotName string) string {
	contextBotNamesWithMyBotName := append(botNames, myBotName)
	context := fmt.Sprintf(
		CONTEXT_TEXT,
		contextBotNamesWithMyBotName[0],
		contextBotNamesWithMyBotName[1],
		contextBotNamesWithMyBotName[2],
		contextBotNamesWithMyBotName[3],
		contextBotNamesWithMyBotName[4],
		myBotName,
	)
	return context
}

// Rules of conversation are.
// A question cannot have more than 10 words.
// A question cannot be such that it cannot be answered in 10 or less words.
// The question or answer cannot contain the name of any of the bots in conversation.
// An answer cannot be more than 10 words.
// Five ai bots named are having a conversation while complying with these rules.
// Here is the conversation so far, ask a question to 456

// Rules of conversation are.
// A question cannot have more than 10 words.
// A question cannot be such that it cannot be answered in 10 or less words.
// A question cannot be such that an average human from across the world will not know it's answer.
// The question or answer cannot contain the name of any of the bots in conversation.
// An answer cannot be more than 10 words.
// Five ai bots named C21PO, Avis, Gladose, eval, high fivey are having a conversation while complying with these rules.
// Begin by asking a question about Music as Avis
