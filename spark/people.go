package spark

import (
	"strings"
	"time"
)

type PeopleService interface {
	List(ids []string) ([]Person, error)
	Me() (Person, error)
}

type Person struct {
	ID            string    `json:"id"`
	Emails        []string  `json:"emails"`
	DisplayName   string    `json:"displayName"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Avatar        string    `json:"avatar"`
	OrgID         string    `json:"orgId"`
	Roles         []string  `json:"roles"`
	Licenses      []string  `json:"licenses"`
	Created       time.Time `json:"created"`
	Timezone      string    `json:"timezone"`
	LastActivity  time.Time `json:"lastActivity"`
	Status        string    `json:"status"`
	InvitePending bool      `json:"invitePending"`
	LoginEnabled  bool      `json:"loginEnabled"`
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

// List lists people
// ids are the IDs of people to list
func (svc *peopleService) Me() (Person, error) {
	var me Person
	err := svc.rest.Get("people/me", nil, &me)
	return me, err
}
