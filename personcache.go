package main

import (
	"github.com/rcorre/spork/spark"
)

type PersonCache interface {
	Get(id string) (string, error)
	Load(ids []string) error
	IsMe(id string) bool
}

type personCache struct {
	svc   spark.PeopleService
	cache map[string]string
	me    string
}

func NewPersonCache(svc spark.PeopleService) (PersonCache, error) {
	me, err := svc.Me()
	if err != nil {
		return nil, err
	}
	return &personCache{
		svc:   svc,
		cache: map[string]string{},
		me:    me.ID,
	}, nil
}

// Get retrieves a single person by id.
// If the person is not cached, an API call is made to look up the id.
// id is the id of the person to fetch
func (p *personCache) Get(id string) (string, error) {
	if name, ok := p.cache[id]; ok {
		return name, nil
	}

	if err := p.Load([]string{id}); err != nil {
		return "", err
	}

	if name, ok := p.cache[id]; ok {
		return name, nil
	}

	// the message may have been sent by a removed user
	return "???", nil
}

// Load makes a batch API request to load multiple people
// ids is the list of person ids to look up
func (p *personCache) Load(ids []string) error {
	list, err := p.svc.List(ids)
	if err != nil {
		return err
	}

	for _, person := range list {
		p.cache[person.ID] = person.DisplayName
	}

	return nil
}

func (p *personCache) IsMe(id string) bool {
	return p.me == id
}
