package store

import "github.com/Vladislav001/golang_study_http_rest_api/internal/app/model"

type UserRepository interface {
	Create(*model.User) error
	FindByEmail(value string) (*model.User, error)
	Find(int2 int) (*model.User, error)
}
