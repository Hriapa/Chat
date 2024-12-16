package store

import (
	"database/sql"
	"webServer/model"
)

type MessageStore struct {
	store *Store
}

func (m *MessageStore) GetRoomID(id1 int, id2 int) (*model.Room, error) {
	var user1, user2 int
	room := &model.Room{}
	if !m.store.Connect {
		return room, ErrorConnectToDatabase
	}
	if id1 < id2 {
		user1 = id1
		user2 = id2
	} else {
		user1 = id2
		user2 = id1
	}
	err := m.store.Db.QueryRow("SELECT id, quantity FROM "+m.store.Config.DbRoomsTable+" WHERE  user_id1 = $1 AND user_id2 = $2", user1, user2).Scan(&room.Id, &room.NumberOfMessage)
	if err == sql.ErrNoRows {
		err = m.store.Db.QueryRow(
			"INSERT INTO "+m.store.Config.DbRoomsTable+" (user_id1, user_id2, quantity) VALUES($1, $2, $3) RETURNING id",
			user1,
			user2,
			0,
		).Scan(&room.Id)
	}
	return room, err
}

func (m *MessageStore) AddMessage(message *model.Message) error {
	if !m.store.Connect {
		return ErrorConnectToDatabase
	}

	_, err := m.store.Db.Exec(
		"INSERT INTO "+m.store.Config.DbMessagesTable+" (index_number, user_id, format, message, reference) VALUES($1, $2, $3, $4, $5)",
		message.IndexNumber,
		message.UserId,
		message.MessageFormat,
		message.Data,
		message.Reference,
	)
	if err != nil {
		return err
	}

	_, err = m.store.Db.Exec("UPDATE "+m.store.Config.DbRoomsTable+" SET quantity = quantity + 1 WHERE id = $1", message.Reference)

	if err != nil {
		return err
	}

	return nil
}

func (m *MessageStore) GetMessages(ref int, maxMessages int, messageNumber int, messages *[]model.Message) error {
	if !m.store.Connect {
		return ErrorConnectToDatabase
	}
	var limit, offset int
	if messageNumber < 20 {
		limit = messageNumber
	} else {
		limit = 20
	}
	offset = maxMessages - messageNumber
	rows, err := m.store.Db.Query("SELECT index_number, user_id, format, message FROM "+m.store.Config.DbMessagesTable+
		" WHERE reference = $1 ORDER BY index_number DESC LIMIT $2 OFFSET $3",
		ref,
		limit,
		offset,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		message := model.Message{}
		err := rows.Scan(&message.IndexNumber, &message.UserId, &message.MessageFormat, &message.Data)
		if err != nil {
			return err
		}
		*messages = append(*messages, message)
	}
	return nil
}

func (m *MessageStore) AddNotReadedMessage(userId int, roomId int) error {
	if !m.store.Connect {
		return ErrorConnectToDatabase
	}
	var res int64
	result, err := m.store.Db.Exec("UPDATE "+m.store.Config.DbNotReadedTable+" SET number = number + 1 WHERE user_id=$1 AND room_id=$2",
		userId,
		roomId,
	)
	if err != nil {
		return err
	}
	res, err = result.RowsAffected()
	if err != nil {
		return err
	}
	if res == 0 {
		_, err = m.store.Db.Exec("INSERT INTO "+m.store.Config.DbNotReadedTable+" (user_id, room_id, number) VALUES($1, $2, $3)",
			userId,
			roomId,
			1,
		)
	}
	return err
}

func (m *MessageStore) DeleteNotReadedMessage(userId int, roomId int) error {
	_, err := m.store.Db.Exec("DELETE FROM "+m.store.Config.DbNotReadedTable+" WHERE user_id=$1 AND room_id=$2",
		userId,
		roomId,
	)
	return err
}

func (m *MessageStore) SelectNotReadMessagesToList(id int, messages map[int]int) error {
	if !m.store.Connect {
		return ErrorConnectToDatabase
	}
	rows, err := m.store.Db.Query("SELECT room_id, number FROM "+m.store.Config.DbNotReadedTable+" WHERE  user_id = $1", id)
	if err != nil {
		return err
	}
	for rows.Next() {
		var (
			room, number int
		)
		err := rows.Scan(&room, &number)
		if err != nil {
			return err
		}
		messages[room] = number
	}
	return nil
}
