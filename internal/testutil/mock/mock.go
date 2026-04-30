package mock

type SetableMock interface {
	Set()
}

type ResetableMock interface {
	Reset()
}

func Set(mocks ...SetableMock) {
	for _, mock := range mocks {
		mock.Set()
	}
}

func Reset(mocks ...ResetableMock) {
	for _, mock := range mocks {
		mock.Reset()
	}
}
