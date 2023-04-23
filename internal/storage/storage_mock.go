package storage

type StorageAccessorMock struct {
	UserRetriever
	GameAccessor
	PlayerAccessor
	MessageCreator
	BotAccessor
	DatabaseTransactionProvider
}

type StorageAccessorMockOption func(*StorageAccessorMock)

func NewStorageAccessorMock(opts ...StorageAccessorMockOption) *StorageAccessorMock {
	mock := &StorageAccessorMock{}
	for _, opt := range opts {
		opt(mock)
	}

	return mock
}

func WithPlayerAccessorMock(mock PlayerAccessor) StorageAccessorMockOption {
	return func(s *StorageAccessorMock) {
		s.PlayerAccessor = mock
	}
}

func WithGameAccessorMock(mock GameAccessor) StorageAccessorMockOption {
	return func(s *StorageAccessorMock) {
		s.GameAccessor = mock
	}
}

func WithMessageCreatorMock(mock MessageCreator) StorageAccessorMockOption {
	return func(s *StorageAccessorMock) {
		s.MessageCreator = mock
	}
}

func WithBotAccessorMock(mock BotAccessor) StorageAccessorMockOption {
	return func(s *StorageAccessorMock) {
		s.BotAccessor = mock
	}
}

func WithDatabaseTransactionProviderMock(mock DatabaseTransactionProvider) StorageAccessorMockOption {
	return func(s *StorageAccessorMock) {
		s.DatabaseTransactionProvider = mock
	}
}
