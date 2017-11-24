package spark

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type PeopleTestSuite struct {
	suite.Suite
}

func TestPeopleTestSuite(t *testing.T) {
	suite.Run(t, new(PeopleTestSuite))
}

func (suite *PeopleTestSuite) TestList() {
	restClient := &RESTClientMock{}
	restClient.On(
		"Get",
		"people",
		map[string]string{
			"id": "one,two,three",
		},
		&struct{ Items []Person }{},
	).Run(func(args mock.Arguments) {
		out := args.Get(2).(*struct{ Items []Person })
		out.Items = []Person{
			{DisplayName: "Foo"},
			{DisplayName: "Bar"},
			{DisplayName: "Baz"},
		}
	}).Return(nil)
	peopleService := NewPeopleService(restClient)

	rooms, err := peopleService.List([]string{"one", "two", "three"})
	suite.Nil(err)
	suite.Equal(rooms, []Person{
		{DisplayName: "Foo"},
		{DisplayName: "Bar"},
		{DisplayName: "Baz"},
	})

	restClient.AssertExpectations(suite.T())
}

func (suite *PeopleTestSuite) TestMe() {
	restClient := &RESTClientMock{}
	restClient.On(
		"Get",
		"people/me",
		map[string]string(nil),
		mock.Anything,
	).Run(func(args mock.Arguments) {
		out := args.Get(2).(*Person)
		*out = Person{
			ID:          "mee-123",
			DisplayName: "Ryan",
		}
	}).Return(nil)
	peopleService := NewPeopleService(restClient)

	me, err := peopleService.Me()
	suite.Nil(err)
	suite.Equal(me.ID, "mee-123")
	suite.Equal(me.DisplayName, "Ryan")

	restClient.AssertExpectations(suite.T())
}
