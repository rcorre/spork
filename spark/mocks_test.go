package spark

import "github.com/stretchr/testify/mock"

// Implementation of common mocks for the spark package

type RESTClientMock struct {
	mock.Mock
}

func (m *RESTClientMock) Get(path string, params map[string]string, out interface{}) error {
	args := m.Called(path, params, out)
	return args.Error(0)
}