package storage

type StorageAccessorMock struct {
	UserRetriever
	InstagramAccountAccessor
}

type StorageAccessorMockOption func(*StorageAccessorMock)

func NewStorageAccessorMock(opts ...StorageAccessorMockOption) *StorageAccessorMock {
	mock := &StorageAccessorMock{}
	for _, opt := range opts {
		opt(mock)
	}

	return mock
}

func WithInstagramAccountAccessorMock(mock InstagramAccountAccessor) StorageAccessorMockOption {
	return func(s *StorageAccessorMock) {
		s.InstagramAccountAccessor = mock
	}
}
