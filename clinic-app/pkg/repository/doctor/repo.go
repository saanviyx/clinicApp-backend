package doctor

import (
	"clinic-app/pkg/repository"
)

type repo struct{}

// New creates a new instance of repository with a database connection
func New() repository.DoctorRepository {
	return &repo{}
}
