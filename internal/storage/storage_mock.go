package storage

type StorageAccessorMock struct {
	UserRetriever
}

type StorageAccessorMockOption func(*StorageAccessorMock)

func NewStorageAccessorMock(opts ...StorageAccessorMockOption) *StorageAccessorMock {
	mock := &StorageAccessorMock{}
	for _, opt := range opts {
		opt(mock)
	}

	return mock
}
