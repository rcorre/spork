package spark

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RoomTestSuite struct {
	suite.Suite
}

func TestRoomTestSuite(t *testing.T) {
	suite.Run(t, new(RoomTestSuite))
}

func (suite *RoomTestSuite) TestList() {
	restClient := &RESTClientMock{}
	restClient.On(
		"Get",
		"rooms",
		map[string]string{},
		&struct{ Items []Room }{},
	).Run(func(args mock.Arguments) {
		out := args.Get(2).(*struct{ Items []Room })
		out.Items = []Room{
			{Title: "Foo"},
			{Title: "Bar"},
			{Title: "Baz"},
		}
	}).Return(nil)
	roomService := NewRoomService(restClient)

	rooms, err := roomService.List()
	suite.Nil(err)
	suite.Equal(rooms, []Room{
		{Title: "Foo"},
		{Title: "Bar"},
		{Title: "Baz"},
	})
}
