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
	svc.On("Me").Return(spark.Person{}, nil)

	svc.On(
		"List",
		[]string{"abc-123"},
	).Return(
		[]spark.Person{{ID: "abc-123", DisplayName: "foo"}},
		nil,
	)

	people, err := NewPersonCache(svc)
	suite.Nil(err)

	name, err := people.Get("abc-123")
	suite.Nil(err)
	suite.Equal("foo", name)

	name, err = people.Get("abc-123")
	suite.Nil(err)
	suite.Equal("foo", name)

	svc.AssertExpectations(suite.T())
}

func (suite *PersonCacheTestSuite) TestIsMe() {
	svc := &mocks.PeopleService{}
	svc.On("Me").Return(spark.Person{
		ID:          "mee-123",
		DisplayName: "Ryan",
	}, nil)

	people, err := NewPersonCache(svc)
	suite.Nil(err)

	suite.True(people.IsMe("mee-123"))
	suite.False(people.IsMe("abc-123"))

	svc.AssertExpectations(suite.T())
}
