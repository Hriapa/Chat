package chat

import (
	"log"
	"time"
	"webServer/model"
	"webServer/protocol"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = pongWait * 9 / 10
	maxMessageSize = 65535
)

type EventType uint8

const (
	Contorl EventType = iota
	Data
)

type Client struct {
	id           int
	conn         *websocket.Conn
	manager      *ChatManager
	control      *protocol.ControlMessage
	data         *protocol.DataMessage
	ack          *protocol.AckMessage
	err          *protocol.ErrorMessage
	message      *model.Message
	assembler    *MessageAssembler
	rooms        map[int]*model.Room // key - room ID (opponent), val - room ID from DB, number of messages in room
	lastMessages []model.Message     // для отправки сообщений из БД
	send         chan []byte
}

func NewClient(id int, conn *websocket.Conn, manager *ChatManager) *Client {
	return &Client{
		id:           id,
		conn:         conn,
		manager:      manager,
		control:      protocol.NewControlMessage(),
		data:         protocol.NewDataMessage(),
		ack:          protocol.NewAckMessage(),
		err:          protocol.NewErrorMessage(),
		message:      model.NewMessage(),
		assembler:    NewMessageAssembler(),
		rooms:        make(map[int]*model.Room),
		lastMessages: make([]model.Message, 0, 20),
		send:         make(chan []byte),
	}
}

func (c *Client) SendUsersList() {
	var (
		i, k int
	)
	nrMessages := make(map[int]int)
	err := c.manager.Store.Message().SelectNotReadMessagesToList(c.id, nrMessages)
	if err != nil {
		log.Println("error with select not readed messages", err)
	}
	c.control.Clear()
	c.control.Type = protocol.UsersListMessageType
	c.control.UsersList.UserListCommand = protocol.SetUserList
	c.control.UsersList.NumberOfElements = len(c.manager.Store.User().UsersList) - 1
	i = 0
	k = 0
	for key, value := range c.manager.Store.User().UsersList {
		if key != c.id {
			c.control.UsersList.List = append(c.control.UsersList.List, protocol.User{Id: key, Name: value.Name, Online: value.Online, NotReadMessages: nrMessages[key]})
			i++
			k++
		}
		if i == 99 || i == c.control.UsersList.NumberOfElements {
			c.send <- c.control.Code()
			i = 0
			c.control.UsersList.UserListCommand = protocol.AddToUserList
			c.control.UsersList.List = c.control.UsersList.List[:0]
		}
	}
}

func (c *Client) ReadMessages() {

	defer c.disconnect()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		var (
			to, roomId int
		)

		messageType, payload, err := c.conn.ReadMessage()

		if err != nil || messageType == websocket.CloseMessage {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}

		if len(payload) == 0 {
			log.Println(`Error. Empty message`)
			continue
		}

		//Control Message
		if payload[0] == uint8(protocol.ControlMessageTitle) {
			protocol.Cleaner(c.control)
			err = protocol.Decoder(c.control, payload)
			if err != nil {
				log.Printf("error decoding control message: %v", err)
				continue
			}
			switch c.control.Type {
			case protocol.MessagesRequestMessageType:
				c.messageRequestProcessing()
			}
			continue
		}

		//Data Message
		if payload[0] == uint8(protocol.DataMessageTitle) {
			protocol.Cleaner(c.data)
			err = protocol.Decoder(c.data, payload)
			if err != nil {
				log.Printf("error decoding message: %v", err)
				continue
			}

			//Если сообщение фрагментировано пробуем его собрать

			if c.data.UserData.Fragmentation.On {
				ok, fullmessage := c.assembler.Processor(c.data)
				// Если сообщение не собрано или ошибка ждём следующего сообщения
				if !ok {
					continue
				}
				c.data.UserData.Data = fullmessage
			}

			// Получаем идентификатор комнаты для  БД
			roomId, err = c.getRoomId(c.data.UserId)

			if err != nil {
				log.Println(`error, save message to DB`, err)
				c.messageDataErrorProcessing()
				continue
			}
			// Сохраняем сообщение в базу данных
			err = c.plaseMessageToDb(roomId)
			if err != nil {
				log.Println(`error, save message to DB`, err)
				c.messageDataErrorProcessing()
				continue
			}

			//Увеличиваем количество сообщений в комнате

			if val, ok := c.rooms[c.data.UserId]; ok {
				val.NumberOfMessage += 1
			}

			c.data.Type = protocol.NewMessage
			// Указываем в отправители текущего пользователя
			if c.data.UserId != 0 {
				to = c.data.UserId

			} else if c.data.Room {
				to = c.data.RoomId
			}
			c.data.UserId = c.id
			c.processingDataMessage(to)
		}

		//Acknowledge message
		if payload[0] == uint8(protocol.AckMessageTitle) {
			protocol.Cleaner(c.ack)
			err = protocol.Decoder(c.ack, payload)
			if err != nil {
				log.Printf("error decoding control message: %v", err)
				continue
			}
			switch c.ack.Type {
			case protocol.Read:
				c.messageReadProcessing()
			}
		}
	}

}

func (c *Client) WriteMessages() {

	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.deleteClient()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed: ", err)
				}
				return
			}
			if err := c.conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
				log.Println(err)
			}
			log.Println("message is sending")

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

// Проверка несобранных фрагментированных сообщений
func (c *Client) AssemblerStoreCleaner() {
	ticker := time.NewTicker(60 * time.Second)

	defer ticker.Stop()

	go func() {
		for {
			<-ticker.C
			c.assembler.assemblerCleane()
		}
	}()
}

// отправка ошибок, если несобранные сообщения удаленны из хранилища

func (c *Client) AssemblingStoreErrorProcessing() {
	for message := range c.assembler.Err {
		protocol.Cleaner(c.err)
		c.err.Type = protocol.DataMessageError
		c.err.IndexNumber = message.Number
		c.err.UserId = message.Id
		c.send <- protocol.Coder(c.err)
	}
}

func (c *Client) deleteClient() {
	c.manager.Delete <- c
}

func (c *Client) plaseMessageToDb(id int) error {
	c.message.Clear()
	c.message.IndexNumber = c.data.IndexNumber
	c.message.UserId = c.id
	c.message.MessageFormat = c.data.Format
	c.message.Reference = id
	c.message.Data = c.data.UserData.Data
	err := c.manager.Store.Message().AddMessage(c.message)
	if err != nil {
		return err
	}
	err = c.manager.Store.Message().AddNotReadedMessage(c.data.UserId, c.id)
	if err != nil {
		log.Println(`error save data to not readed messages table`, err)
	}
	return nil
}
