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
	ID            string     `json:"id,omitempty"`
	Emails        []string   `json:"emails,omitempty"`
	DisplayName   string     `json:"displayName,omitempty"`
	FirstName     string     `json:"firstName,omitempty"`
	LastName      string     `json:"lastName,omitempty"`
	Avatar        string     `json:"avatar,omitempty"`
	OrgID         string     `json:"orgId,omitempty"`
	Roles         []string   `json:"roles,omitempty"`
	Licenses      []string   `json:"licenses,omitempty"`
	Created       *time.Time `json:"created,omitempty"`
	Timezone      string     `json:"timezone,omitempty"`
	LastActivity  *time.Time `json:"lastActivity,omitempty"`
	Status        string     `json:"status,omitempty"`
	InvitePending bool       `json:"invitePending,omitempty"`
	LoginEnabled  bool       `json:"loginEnabled,omitempty"`
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
