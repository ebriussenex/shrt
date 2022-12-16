package shorten

import (
	"context"
	"fmt"

	"github.com/ebriussenex/shrt/internal/model"
	"github.com/google/uuid"
)

type Storage interface {
	Put(ctx context.Context, shrt model.Shortened) (*model.Shortened, error)
	Get(ctx context.Context, id string) (*model.Shortened, error)
	IncrVisits(ctx context.Context, id string) error
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) Shorten(ctx context.Context, req model.ShortenedReq) (*model.Shortened, error) {
	id := req.Id
	if id == "" {
		id = Shorten(uuid.New().ID())
	}

	inShrt := model.Shortened{
		Id:        id,
		Url:       req.Url,
		CreatedBy: req.CreatedBy,
	}

	shrt, err := s.storage.Put(ctx, inShrt)
	if err != nil {
		return nil, err
	}

	return shrt, nil
}

func (s *Service) Get(ctx context.Context, id string) (*model.Shortened, error) {
	shortening, err := s.storage.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return shortening, nil
}

func (s *Service) Redirect(ctx context.Context, id string) (string, error) {
	shortened, err := s.storage.Get(ctx, id)
	if err != nil {
		return "", err
	}

	if err := s.storage.IncrVisits(ctx, id); err != nil {
		return "", fmt.Errorf("failed to redirect for id: %s, err: %w", id, err)
	}
	
	return shortened.Url, nil
}

