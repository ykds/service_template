package service

import (
	"github.com/redis/go-redis/v9"
	"service_template/internal/repository"
)

func NewService(repo *repository.Repository, cache *redis.Client) *Service {
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
