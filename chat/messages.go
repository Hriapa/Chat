package chat

import (
	"log"
	"webServer/model"
	"webServer/protocol"
)

func (c *Client) broadcastMessage(message []byte) {
	for client := range c.manager.clients {
		if c.conn != client.conn {
			client.send <- message
		}
	}
}

func (c *Client) connect() {

	protocol.Cleaner(c.control)

	c.control.Type = protocol.ConnectMessageType
	c.control.UserId = c.id

	c.broadcastMessage(protocol.Coder(c.control))
}

func (c *Client) disconnect() {

	protocol.Cleaner(c.control)

	c.control.Type = protocol.DisconnectMessageType
	c.control.UserId = c.id

	c.broadcastMessage(protocol.Coder(c.control))
}

func (c *ChatManager) registerNewUser(user *model.UserName) {

	control := protocol.NewControlMessage()

	control.Type = protocol.RegistrationMessageType
	control.UserId = user.Id
	control.UserName = user.Name

	message := protocol.Coder(control)

	for client := range c.clients {
		client.send <- message
	}
}

func (c *Client) messageRequestProcessing() {
	_, err := c.getRoomId(c.control.RoomId)

	if err != nil {
		log.Println(`error, get id from DB`, err)
		return
	}
	if c.chekAndSendMessages() {
		protocol.Cleaner(c.data)
		for _, val := range c.lastMessages {
			c.data.Type = protocol.OldMessage
			c.data.Format = val.MessageFormat
			c.data.IndexNumber = val.IndexNumber
			c.data.UserId = val.UserId
			c.data.Room = true
			c.data.RoomId = c.control.RoomId
			c.data.UserData.Data = val.Data
			c.send <- protocol.Coder(c.data)
		}
	}
}

func (c *Client) messageReadProcessing() {
	var err error
	if c.ack.UserId != 0 {
		err = c.manager.Store.Message().DeleteNotReadedMessage(c.id, c.ack.UserId)
		if err != nil {
			log.Println(`error delete not readed message`, err)
			return
		}
	}
}

func (c *Client) messageDataErrorProcessing() {
	protocol.Cleaner(c.err)
	c.err.Type = protocol.DataMessageError
	c.err.IndexNumber = c.data.IndexNumber
	if c.data.Room {
		c.err.RoomId = c.data.RoomId
	} else {
		c.err.UserId = c.data.UserId
	}
	c.send <- protocol.Coder(c.err)
}

// Получаем идентификатор комнаты для  БД
func (c *Client) getRoomId(id int) (int, error) {
	if val, ok := c.rooms[id]; ok {
		return val.Id, nil
	}

	room, err := c.manager.Store.Message().GetRoomID(c.id, id)
	if err != nil {
		return 0, err
	}

	c.rooms[id] = room

	return room.Id, nil
}

// Подготовак сообщений из БД для отправки
func (c *Client) chekAndSendMessages() bool {
	lastMessageNumber := c.control.MessagesRequest.LastMessageNumber
	numberOfMessages := c.rooms[c.control.RoomId].NumberOfMessage
	if lastMessageNumber == 1 {
		return false
	}
	if lastMessageNumber == 0 { // Отправляем количество сообщений
		c.control.Type = protocol.NumberOfMessagesMessageType
		c.control.NumberOfMessages = c.rooms[c.control.RoomId].NumberOfMessage
		c.send <- c.control.Code()
		lastMessageNumber = c.rooms[c.control.RoomId].NumberOfMessage + 1
		if lastMessageNumber == 1 {
			return false
		}
	}
	ref := c.rooms[c.control.RoomId].Id
	c.lastMessages = c.lastMessages[:0]
	err := c.manager.Store.Message().GetMessages(ref, numberOfMessages, lastMessageNumber-1, &c.lastMessages)
	if err != nil {
		log.Println(`error, select messages`)
		return false
	}
	return true
}
