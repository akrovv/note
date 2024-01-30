package service

import (
	"context"
	"note/internal/models"
)

type Storage interface {
	Get(context.Context, string) (*models.Note, error)
	Create(context.Context, string) error
	Update(context.Context, string, string) error
	Delete(context.Context, string) error
	GetAll(context.Context, string) ([]*models.Note, error)
}

type service struct {
	storage Storage
}

func NewService(storage Storage) *service {
	return &service{storage: storage}
}

func (s *service) Create(ctx context.Context, dto CreateNote) error {
	return s.storage.Create(ctx, dto.Text)
}

func (s *service) Update(ctx context.Context, dto UpdateNote) error {
	return s.storage.Update(ctx, dto.ID, dto.Text)
}

func (s *service) Delete(ctx context.Context, dto DeleteNote) error {
	return s.storage.Delete(ctx, dto.ID)
}

func (s *service) Get(ctx context.Context, dto GetNote) (*models.Note, error) {
	return s.storage.Get(ctx, dto.ID)
}

func (s *service) GetAll(ctx context.Context, dto GetNotes) ([]*models.Note, error) {
	return s.storage.GetAll(ctx, dto.OrderBy)
}
