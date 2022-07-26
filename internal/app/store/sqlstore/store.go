package sqlstore

import (
	"database/sql"
	"github.com/Vladislav001/golang_study_http_rest_api/internal/app/store"

	_ "github.com/lib/pq" // анонимный импорт, чтобы методы не импортились
)

type Store struct {
	db             *sql.DB
	userRepository *UserRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// eg. store.User().Create()
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}
	return s.userRepository
}
