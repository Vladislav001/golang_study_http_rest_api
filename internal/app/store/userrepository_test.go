package store_test

import (
	"github.com/Vladislav001/golang_study_http_rest_api/internal/app/model"
	"github.com/Vladislav001/golang_study_http_rest_api/internal/app/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	s, teardown := store.TestStore(t, databaseURL)
	defer teardown("users")

	u, err := s.User().Create(&model.User{
		Email: "test@mail.ru",
	})

	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	s, teardown := store.TestStore(t, databaseURL)
	defer teardown("users")

	email := "test@mail.ru"
	_, err := s.User().FindByEmail(email)
	assert.Error(t, err)

	_, err = s.User().Create(&model.User{
		Email: "test@mail.ru",
	})
	u, err := s.User().FindByEmail(email)

	assert.NoError(t, err)
	assert.NotNil(t, u)
}