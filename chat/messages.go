package chat

import (
	"log"
	"time"
	"webServer/model"
	"webServer/protocol"
)

func (c *Client) broadcastMessage(message []byte) {
	for _, client := range c.manager.clients {
		if c.conn != client.conn {
			client.send <- message
		}
	}
}

func (c *Client) singleClientMessage(id int, message []byte) {
	val, ok := c.manager.clients[id]
	if !ok {
		log.Printf(`client is offline`)
		return
	}
	val.send <- message
}

func (c *Client) sendDataMessage(id int, message []byte) {
	if id == 0 {
		c.broadcastMessage(message)
	} else {
		c.singleClientMessage(id, message)
	}
}

// Data Message

func (c *Client) processingDataMessage(id int) {
	if len(c.data.UserData.Data) > maxFragmentLength {
		c.data.UserData.Fragmentation.On = true
		result := MessageFragmentation(c.data.UserData.Data, maxFragmentLength)
		for _, val := range result {
			if val.SequenceNumber == 1 {
				c.data.UserData.Fragmentation.FragmentType = protocol.FirstFragment
			} else if val.SequenceNumber == uint8(len(result)) {
				c.data.UserData.Fragmentation.FragmentType = protocol.LastFragment
			} else {
				c.data.UserData.Fragmentation.FragmentType = protocol.MiidleFragment
			}
			c.data.UserData.Fragmentation.Counter = val.SequenceNumber
			c.data.UserData.Data = val.Data
			c.sendDataMessage(id, protocol.Coder(c.data))
		}
		return
	}
	c.data.UserData.Fragmentation.On = false
	c.sendDataMessage(id, protocol.Coder(c.data))
}

// Connect Message Processing

func (c *Client) connect() {

	protocol.Cleaner(c.control)

	c.control.Type = protocol.ConnectMessageType
	c.control.UserId = c.id

	c.broadcastMessage(protocol.Coder(c.control))
}

// Disconnect Message Processing

func (c *Client) disconnect() {

	protocol.Cleaner(c.control)

	c.control.Type = protocol.DisconnectMessageType
	c.control.UserId = c.id

	c.broadcastMessage(protocol.Coder(c.control))
}

// Registration Message Processing

func (c *ChatManager) registerNewUser(user *model.UserName) {

	control := protocol.NewControlMessage()

	control.Type = protocol.RegistrationMessageType
	control.UserId = user.Id
	control.UserName = user.Name

	message := protocol.Coder(control)

	for _, client := range c.clients {
		client.send <- message
	}
}

// User Info Message Processing

func (c *Client) userInfoRequestProcessing() {
	user := &model.UserInfo{
		Id: c.id,
	}
	err := c.manager.Store.User().GetUserInfo(user)
	if err != nil {
		log.Println("error select user info from Db:", err.Error())
		c.userInfoErrorProcessing()
		return
	}
	cleaner(c.control)
	c.control.Type = protocol.UserInfoResponseMessageType
	c.control.UserInfo = &protocol.UserInfo{
		Name:       user.Name,
		Surname:    user.Surname,
		FamilyName: user.Familyname,
		BirthDate:  dateConverter(user.Birthdate),
	}
	c.send <- coder(c.control)
}

func (c *Client) userInfoUpdateProcessing() {

	date := time.Date(int(c.control.UserInfo.BirthDate.Year), time.Month(c.control.UserInfo.BirthDate.Month), int(c.control.UserInfo.BirthDate.Day), 0, 0, 0, 0, time.UTC)

	user := &model.UserInfo{
		Id:         c.id,
		Name:       c.control.UserInfo.Name,
		Familyname: c.control.UserInfo.FamilyName,
		Surname:    c.control.UserInfo.Surname,
		Birthdate:  date,
	}

	err := c.manager.Store.User().UpdateUserInfo(user)
	if err != nil {
		log.Println("error update user info:", err.Error())
		c.userInfoErrorProcessing()
		return
	}
	// TO DO: SEND OK
}

// Message Request Processing

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
			c.processingDataMessage(c.id)
		}
	}
}

// Ack Processing

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

// Error procrssing

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

func (c *Client) userInfoErrorProcessing() {
	cleaner(c.err)
	c.err.Type = protocol.ControlMessageError
	c.err.ControlMessageType = protocol.UserInfoResponseMessageType
	c.send <- coder(c.err)
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

// Convert date from string to protocol format

func dateConverter(date time.Time) protocol.BirthDate {
	if date.Year() == -1 {
		return protocol.BirthDate{
			Year:  0,
			Month: 0,
			Day:   0,
		}
	}
	birthDate := protocol.BirthDate{
		Year:  uint16(date.Year()),
		Month: uint8(date.Month()),
		Day:   uint8(date.Day()),
	}
	return birthDate
}
