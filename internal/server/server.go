package server

import (
	"github.com/vipulvpatil/airetreat-go/internal/clients/openai"
	"github.com/vipulvpatil/airetreat-go/internal/config"
	"github.com/vipulvpatil/airetreat-go/internal/storage"
	pb "github.com/vipulvpatil/airetreat-go/protos"
)

type AiRetreatGoService struct {
	pb.UnsafeAiRetreatGoServer
	storage      storage.StorageAccessor
	openAiClient openai.Client
	config       *config.Config
}

type ServerDependencies struct {
	Storage      storage.StorageAccessor
	OpenAiClient openai.Client
	Config       *config.Config
}

func NewServer(deps ServerDependencies) (*AiRetreatGoService, error) {
	return &AiRetreatGoService{
		storage:      deps.Storage,
		openAiClient: deps.OpenAiClient,
		config:       deps.Config,
	}, nil
}
