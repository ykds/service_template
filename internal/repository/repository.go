package repository

import (
	"errors"
	"service_template/pkg/db"
)

var (
	RecordNotFound = errors.New("record not found")
)

func NewRepository(db *db.DB) *Repository {
	return &Repository{
		ExampleRepository: NewExampleRepository(db),
	}
}

type Repository struct {
	ExampleRepository ExampleRepository
}

type ExampleRepository interface {
}

func NewExampleRepository(db *db.DB) ExampleRepository {
	return &exampleZoneRepository{
		db: db,
	}
}

type exampleZoneRepository struct {
	db *db.DB
}
