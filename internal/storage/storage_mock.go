package storage

type StorageAccessorMock struct {
	UserRetriever
	GameAccessor
	PlayerCreator
	MessageCreator
	TxProvider
}

type StorageAccessorMockOption func(*StorageAccessorMock)

func NewStorageAccessorMock(opts ...StorageAccessorMockOption) *StorageAccessorMock {
	mock := &StorageAccessorMock{}
	for _, opt := range opts {
		opt(mock)
	}

	return mock
}

func WithPlayerCreatorMock(mock PlayerCreator) StorageAccessorMockOption {
	return func(s *StorageAccessorMock) {
		s.PlayerCreator = mock
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
