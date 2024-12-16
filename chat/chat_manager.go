package chat

import (
	"log"
	"sync"
	"webServer/model"
	"webServer/store"

	"github.com/gorilla/websocket"
)

const (
	wsReadBuffer      = 1024
	wsWriteBuffer     = 1024
	maxFragmentLength = 32767
)

type ChatManager struct {
	sync.RWMutex
	Upgrader  websocket.Upgrader
	clients   map[int]*Client
	NewUser   chan *model.UserName
	Register  chan *Client
	Delete    chan *Client
	RoomsList map[int][]int
	Store     *store.Store
}

func NewChatManager(s *store.Store) *ChatManager {
	return &ChatManager{
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  wsReadBuffer,
			WriteBufferSize: wsWriteBuffer,
		},
		clients:   make(map[int]*Client),
		NewUser:   make(chan *model.UserName),
		Register:  make(chan *Client),
		Delete:    make(chan *Client),
		RoomsList: make(map[int][]int),
		Store:     s,
	}
}

func (c *ChatManager) createNewUser(user *model.UserName) {
	c.Lock()
	defer c.Unlock()

	userParametrs := &model.UserInList{Name: user.Name, Online: user.Online}
	c.Store.User().UsersList[user.Id] = userParametrs
	c.registerNewUser(user)
}

func (c *ChatManager) registerClient(client *Client) {
	c.Lock()
	defer c.Unlock()

	c.clients[client.id] = client
	if val, ok := c.Store.User().UsersList[client.id]; ok {
		val.Online = true
	}
	client.connect()

	log.Println("new client register")
}

func (c *ChatManager) deleteClient(client *Client) error {
	c.Lock()
	defer c.Unlock()

	var err error
	if _, ok := c.clients[client.id]; ok {
		client.disconnect()
		if val, ok := c.Store.User().UsersList[client.id]; ok {
			val.Online = false
		}
		err = client.conn.Close()
		delete(c.clients, client.id)
		log.Println("client deleted")
	}
	return err
}

func (c *ChatManager) Run() {
	for {
		select {
		case user := <-c.NewUser:
			c.createNewUser(user)
		case client := <-c.Register:
			c.registerClient(client)
		case client := <-c.Delete:
			c.deleteClient(client)
		}
	}
}
