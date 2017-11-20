package main

import (
	"testing"

	"github.com/rcorre/spork/mocks"
	"github.com/rcorre/spork/spark"
	"github.com/stretchr/testify/suite"
)

type PersonCacheTestSuite struct {
	suite.Suite
}

func TestPersonCacheTestSuite(t *testing.T) {
	suite.Run(t, new(PersonCacheTestSuite))
}

func (suite *PersonCacheTestSuite) TestGet() {
	svc := &mocks.PeopleService{}
	svc.On(
		"List",
		[]string{"abc-123"},
	).Return(
		[]spark.Person{{ID: "abc-123", DisplayName: "foo"}},
		nil,
	)

	people := NewPersonCache(svc)

	name, err := people.Get("abc-123")
	suite.Nil(err)
	suite.Equal("foo", name)

	name, err = people.Get("abc-123")
	suite.Nil(err)
	suite.Equal("foo", name)

	svc.AssertExpectations(suite.T())
}
