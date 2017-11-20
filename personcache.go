package main

import (
	"fmt"

	"github.com/rcorre/spork/spark"
)

type PersonCache interface {
	Get(id string) (string, error)
	Load(ids []string) error
}

type personCache struct {
	svc   spark.PeopleService
	cache map[string]string
}

func NewPersonCache(svc spark.PeopleService) PersonCache {
	return &personCache{
		svc:   svc,
		cache: map[string]string{},
	}
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

	return "", fmt.Errorf("Could not find person %q", id)
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
