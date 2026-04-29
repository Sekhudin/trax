package mock

type Mock interface {
	Reset()
}

func Reset(mocks ...Mock) {
	for _, mock := range mocks {
		mock.Reset()
	}
}
