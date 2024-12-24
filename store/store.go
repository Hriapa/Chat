package store

import (
	"database/sql"
	"fmt"
	"webServer/model"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Store struct {
	Db           *sql.DB
	Connect      bool
	Config       *Config
	userStore    *UsersStore
	messageStore *MessageStore
}

func NewDataBase() *Store {
	return &Store{
		Db:           nil,
		Connect:      false,
		Config:       NewConfig(),
		userStore:    nil,
		messageStore: nil,
	}
}

func (s *Store) Open() error {

	configStr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		s.Config.DbHost, s.Config.DbPort, s.Config.DbUser, s.Config.DbPassword, s.Config.DbName)

	db, err := sql.Open("postgres", configStr)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	// Migrations

	if err = migration(db); err != nil {
		return err
	}

	s.Db = db
	s.Connect = true

	return nil
}

func (s *Store) Close() {
	s.Db.Close()
}

// migration

func migration(db *sql.DB) error {

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://store/migrations",
		"postgres", driver,
	)
	if err != nil {
		return err
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

//Инициализация хранилища пользователей

func (s *Store) User() *UsersStore {
	if s.userStore == nil {
		s.userStore = &UsersStore{
			store:     s,
			UsersList: make(map[int]*model.UserInList),
		}
	}

	return s.userStore
}

//Инициализация хранилища сообщений

func (s *Store) Message() *MessageStore {
	if s.messageStore == nil {
		s.messageStore = &MessageStore{
			store: s,
		}
	}

	return s.messageStore
}
