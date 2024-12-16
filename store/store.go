package store

import (
	"database/sql"
	"fmt"
	"webServer/model"

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

	// Создаём таблицы, если они отсутствуют

	s.Db = db
	s.Connect = true
	return nil
}

func (s *Store) Close() {
	s.Db.Close()
}

// Создание таблиц

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
