package store

// *На win10 надо ставить постгрес + руками создать БД
//(https://winitpro.ru/index.php/2019/10/25/ustanovka-nastrojka-postgresql-v-windows/)

import (
	"database/sql"

	_ "github.com/lib/pq" // анонимный импорт, чтобы методы не импортились
)

type Store struct {
	config *Config
	db     *sql.DB
}

func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

func (s *Store) Open() error {
	db, err := sql.Open("postgres", s.config.DatabaseURL)
	if err != nil {
		return err
	}

	// т.к ленивое подключение к БД на самом деле => точно проверим сами, что ок все
	if err := db.Ping(); err != nil {
		return err
	}

	s.db = db

	return nil
}

func (s *Store) Close() {
	s.db.Close()
}
