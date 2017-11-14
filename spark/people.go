package spark

import "strings"

type PeopleService interface {
	List(ids []string) ([]Person, error)
}

type Person struct {
	ID            string
	Emails        []string
	DisplayName   string
	FirstName     string
	LastName      string
	Avatar        string
	OrgID         string
	Roles         []string
	Licenses      []string
	Created       string
	Timezone      string
	LastActivity  string
	Status        string
	InvitePending bool
	LoginEnabled  bool
}

type peopleService struct {
	rest RESTClient
}

func NewPeopleService(rest RESTClient) PeopleService {
	return &peopleService{rest: rest}
}

// List lists people
// ids are the IDs of people to list
func (svc *peopleService) List(ids []string) ([]Person, error) {
	var list struct {
		Items []Person
	}
	params := map[string]string{
		"id": strings.Join(ids, ","),
	}
	err := svc.rest.Get("people", params, &list)
	return list.Items, err
}
