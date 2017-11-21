package mocks

import (
	"github.com/rcorre/spork/spark"
	"github.com/stretchr/testify/mock"
)

// Implementation of common mocks for the spark package

type PeopleService struct {
	mock.Mock
}

func (m *PeopleService) List(ids []string) ([]spark.Person, error) {
	args := m.Called(ids)
	return args.Get(0).([]spark.Person), args.Error(1)
}

type MessageService struct {
	mock.Mock
}

func (m *MessageService) List(id string) ([]spark.Message, error) {
	args := m.Called(id)
	return args.Get(0).([]spark.Message), args.Error(1)
}

type PersonCache struct {
	mock.Mock
}

func (m *PersonCache) Get(id string) (string, error) {
	args := m.Called(id)
	return args.String(0), args.Error(1)
}

func (m *PersonCache) Load(ids []string) error {
	args := m.Called(ids)
	return args.Error(0)
}
