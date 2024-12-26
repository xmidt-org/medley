package medley

import "github.com/stretchr/testify/mock"

type MockLocator[S Service] struct {
	mock.Mock
}

func (m *MockLocator[S]) Find(object []byte) (S, error) {
	args := m.Called(object)

	svc, _ := args.Get(0).(S)
	return svc, args.Error(1)
}

func (m *MockLocator[S]) ExpectFindSuccess(object any, result S) *mock.Call {
	return m.On("Find", object).Return(result, error(nil))
}

func (m *MockLocator[S]) ExpectFindFail(object any, err error) *mock.Call {
	var zero S
	return m.On("Find", object).Return(zero, err)
}

func (m *MockLocator[S]) ExpectFindNoServices(object any) *mock.Call {
	return m.ExpectFindFail(object, ErrNoServices)
}
