package model

import (
	"webServer/protocol"
)

type Direction uint8

const (
	Incoming Direction = iota + 1
	Outgoing
)

type Message struct {
	IndexNumber   int                        // порядковый номер
	UserId        int                        // id автора
	MessageFormat protocol.DataMessageFormat // Тип сообщения (текст, картинка)
	Data          []byte                     // Сообщение
	Reference     int                        // ссылка на комнату
	//Date
}

func NewMessage() *Message {
	return &Message{}
}

func (m *Message) Clear() {
	m.IndexNumber = 0
	m.UserId = 0
	m.MessageFormat = 0
	m.Data = nil
}

type Room struct {
	Id              int
	NumberOfMessage int
}
