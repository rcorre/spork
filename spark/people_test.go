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
			Person{DisplayName: "Foo"},
			Person{DisplayName: "Bar"},
			Person{DisplayName: "Baz"},
		}
	}).Return(nil)
	peopleService := NewPeopleService(restClient)

	rooms, err := peopleService.List([]string{"one", "two", "three"})
	suite.Nil(err)
	suite.Equal(rooms, []Person{
		Person{DisplayName: "Foo"},
		Person{DisplayName: "Bar"},
		Person{DisplayName: "Baz"},
	})
}
