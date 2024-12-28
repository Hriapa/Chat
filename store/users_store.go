package store

import (
	"time"
	"webServer/model"
)

type UsersStore struct {
	store     *Store
	UsersList map[int]*model.UserInList
}

func (u *UsersStore) AddUser(user *model.User) error {
	if !u.store.Connect {
		return ErrorConnectToDatabase
	}

	if err := user.EncryptPassword(); err != nil {
		return err
	}

	if err := u.store.Db.QueryRow(
		"INSERT INTO "+u.store.Config.DbUsersTable+" (login, encrypted_password, name, familyname, surname, birthdate) VALUES($1, $2, $3, $4, $5, $6) RETURNING id",
		user.Login,
		user.EncryptedPassword,
		``,
		``,
		``,
		time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC),
	).Scan(&user.Id); err != nil {
		return err
	}
	return nil
}

func (u *UsersStore) FindByLogin(login string) (*model.User, error) {
	if !u.store.Connect {
		return nil, ErrorConnectToDatabase
	}
	user := &model.User{}

	row := u.store.Db.QueryRow("SELECT id, login, encrypted_password from "+u.store.Config.DbUsersTable+" where login = $1", login)
	if err := row.Scan(&user.Id, &user.Login, &user.EncryptedPassword); err != nil {
		return nil, ErrorLoginNotFound
	}
	return user, nil
}

func (u *UsersStore) FindById(id int) (*model.User, error) {
	if !u.store.Connect {
		return nil, ErrorConnectToDatabase
	}
	user := &model.User{}

	row := u.store.Db.QueryRow("SELECT id, login, encrypted_password from "+u.store.Config.DbUsersTable+" where id = $1", id)
	if err := row.Scan(&user.Id, &user.Login, &user.EncryptedPassword); err != nil {
		return nil, ErrorLoginNotFound
	}
	return user, nil
}

func (u *UsersStore) GetUserInfo(user *model.UserInfo) error {
	if !u.store.Connect {
		return ErrorConnectToDatabase
	}
	row := u.store.Db.QueryRow("SELECT login, name, familyname, surname, birthdate from "+u.store.Config.DbUsersTable+" where id = $1", user.Id)
	if err := row.Scan(&user.Login, &user.Name, &user.Familyname, &user.Surname, &user.Birthdate); err != nil {
		return ErrorUserNotFound
	}
	return nil
}

func (u *UsersStore) UpdateUserInfo(user *model.UserInfo) error {
	if !u.store.Connect {
		return ErrorConnectToDatabase
	}
	_, err := u.store.Db.Exec("update "+u.store.Config.DbUsersTable+" set name = $1, familyname = $2, surname = $3, birthdate = $4 where id = $5",
		user.Name, user.Familyname, user.Surname, user.Birthdate, user.Id)
	if err != nil {
		return err
	}
	return nil
}

func (u *UsersStore) UpdatePassword(user *model.User) error {
	if !u.store.Connect {
		return ErrorConnectToDatabase
	}
	_, err := u.store.Db.Exec("update "+u.store.Config.DbUsersTable+" set encrypted_password = $1 where id = $2",
		user.EncryptedPassword, user.Id)
	if err != nil {
		return err
	}
	return nil
}

func (u *UsersStore) GetUsersList() error {
	if !u.store.Connect {
		return ErrorConnectToDatabase
	}

	var (
		fullName string
		id       int
	)

	rows, err := u.store.Db.Query("select id, login, name, familyname, surname from " + u.store.Config.DbUsersTable)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		fullName = ""
		var (
			name, familyname, surname, login string
		)
		err := rows.Scan(&id, &login, &name, &familyname, &surname)
		if err != nil {
			return err
		}
		if name != "" {
			fullName = name
		}
		if surname != "" && fullName != "" {
			fullName = fullName + " " + surname
		}
		if familyname != "" && fullName != "" {
			fullName = fullName + " " + familyname
		}
		if fullName == "" {
			fullName = login
		}
		user := &model.UserInList{Name: fullName, Online: false}
		u.UsersList[id] = user
	}
	return nil
}
