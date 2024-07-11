package service

import (
	"service_template/internal/repository"
	"service_template/pkg/cache"
)

func NewService(repo *repository.Repository, cache cache.Cache) *Service {
	return &Service{
		ExampleService: NewExampleService(repo),
	}
}

type Service struct {
	ExampleService ExampleService
}

type ExampleService interface {
}

func NewExampleService(repo repository.ExampleRepository) ExampleService {
	return &exampleService{
		repo: repo,
	}
}

type exampleService struct {
	repo repository.ExampleRepository
}
