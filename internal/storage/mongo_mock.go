package shorten_store

import (
	"context"
	"sync"
	"time"

	"github.com/ebriussenex/shrt/internal/model"
)

type mongoStoreMock struct {
	m sync.Map
}

func NewMongoDbMock() *mongoStoreMock {
	return &mongoStoreMock{}
}

func (s *mongoStoreMock) Put(_ context.Context, shortened model.Shortened) (*model.Shortened, error) {
	if _, exists := s.m.Load(shortened.Id); exists {
		return nil, model.ErrIdAlreadyExists
	}
	shortened.CreatedAt = time.Now().UTC()
	s.m.Store(shortened.Id, shortened)
	return &shortened, nil
}

func (s *mongoStoreMock) Get(_ context.Context, identifier string) (*model.Shortened, error) {
	v, ok := s.m.Load(identifier)
	if !ok {
		return nil, model.ErrNotFound
	}
	shortened := v.(model.Shortened)
	return &shortened, nil
}

func (s *mongoStoreMock) IncrVisits(_ context.Context, identifier string) error {
	v, ok := s.m.Load(identifier)
	if !ok {
		return model.ErrNotFound
	}
	shortened := v.(model.Shortened)
	shortened.VisitCount++
	shortened.UpdatedAt = time.Now()
	s.m.Store(identifier, shortened)
	return nil
}
