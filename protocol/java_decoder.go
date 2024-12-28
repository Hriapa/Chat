package protocol

import (
	"log"
	"os"
)

// Control Message

func (c *ControlMessage) ConnectTestMessage() {
	c.Clear()
	c.Type = ConnectMessageType
	c.UserId = 10
}

func (c *ControlMessage) DisConnectTestMessage() {
	c.Clear()
	c.Type = DisconnectMessageType
	c.UserId = 1024
}

func (c *ControlMessage) RegistrationTestMessage() {
	c.Clear()
	c.Type = RegistrationMessageType
	c.UserId = 66000
	c.UserName = `Stepan`
}

func (c *ControlMessage) UserUpdateTestMessage() {
	c.Clear()
	c.Type = UserUpdateMessageType
	c.UserId = 1024
	c.UserName = `Ivan`
}

func (c *ControlMessage) UserInfoTestMessage() {
	c.Clear()
	c.Type = UserInfoResponseMessageType
	c.UserInfo.Name = `Ivan`
	c.UserInfo.Surname = `Ivanovich`
	c.UserInfo.FamilyName = `Ivanov`
	c.UserInfo.BirthDate.Day = 1
	c.UserInfo.BirthDate.Month = 10
	c.UserInfo.BirthDate.Year = 2000
}

func (c *ControlMessage) UsersListTestMessage() {
	c.Clear()
	c.Type = UsersListMessageType
	c.UsersList.UserListCommand = SetUserList
	c.UsersList.NumberOfElements = 3
	c.UsersList.List = []User{{Id: 1, Name: "Petr", Online: false, NotReadMessages: 1},
		{Id: 1024, Name: "Ivan", Online: true, NotReadMessages: 4},
		{Id: 10, Name: "Olga", Online: true, NotReadMessages: 0}}
}

func (c *ControlMessage) ControlCommandTestCoder(command uint8) []byte {
	switch ControlMessageType(command) {
	case ConnectMessageType:
		c.ConnectTestMessage()
	case DisconnectMessageType:
		c.DisConnectTestMessage()
	case RegistrationMessageType:
		c.RegistrationTestMessage()
	case UserUpdateMessageType:
		c.UserUpdateTestMessage()
	case UserInfoResponseMessageType:
		c.UserInfoTestMessage()
	case UsersListMessageType:
		c.UsersListTestMessage()
	default:
		log.Println(`Unknown Control Message Type`)
		return nil
	}

	return Coder(c)
}

//Data Message

func (d *DataMessage) SmallTextMessage() {
	d.Clear()
	d.Type = NewMessage
	d.Format = Text
	d.IndexNumber = 0
	d.UserId = 10
	d.Room = true
	d.RoomId = 1024
	d.UserData.Data = []byte(`Короткое сообщение`)
}

func (d *DataMessage) BigTextMessage() {
	d.Clear()
	d.Type = OldMessage
	d.Format = Text
	d.IndexNumber = 1024
	d.UserId = 10
	d.UserData.Data = []byte(`Если позволять себе шутить, люди не воспринимают тебя всерьёз. И эти самые люди не понимают, что есть многое, чего нельзя выдержать, если не шутить.`)
}

func (d *DataMessage) FirstFragment() {
	d.Clear()
	d.Type = UpdateMessage
	d.Format = Text
	d.IndexNumber = 250
	d.UserId = 10
	d.UserData.Fragmentation.On = true
	d.UserData.Fragmentation.FragmentType = FirstFragment
	d.UserData.Fragmentation.Counter = 0
	d.UserData.Data = []byte(`First fragment`)
}

func (d *DataMessage) LastFragment() {
	d.Clear()
	d.Type = UpdateMessage
	d.Format = Text
	d.IndexNumber = 10
	d.UserId = 10
	d.UserData.Fragmentation.On = true
	d.UserData.Fragmentation.FragmentType = LastFragment
	d.UserData.Fragmentation.Counter = 10
	d.UserData.Data = []byte(`Last fragment`)
}

func (d *DataMessage) Image() {
	const path = `image\test.jpg`
	d.Clear()
	d.Type = NewMessage
	d.Format = Image
	d.IndexNumber = 10
	d.UserId = 10
	d.UserData.Fragmentation.On = false
	d.UserData.Fragmentation.FragmentType = 0
	d.UserData.Fragmentation.Counter = 0
	d.UserData.Data = readImege(path)
}

func (d *DataMessage) DataMessageTestCoder(command uint8) []byte {
	switch command {
	case 1:
		d.SmallTextMessage()
	case 2:
		d.BigTextMessage()
	case 3:
		d.FirstFragment()
	case 4:
		d.LastFragment()
	case 5:
		d.Image()
	default:
		log.Println(`Unknown test for Data Message`)
		return nil
	}

	return Coder(d)
}

func readImege(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		log.Println(`no file`, err)
		return nil
	}
	out := make([]byte, 32000)
	n, err := file.Read(out)
	if err != nil {
		log.Println(err)
		return nil
	}
	return out[:n]
}

// Acknowledge Message

func (a *AckMessage) ControlUserUpdate() {
	a.Clear()
	a.Type = ControlAck
	a.ControlAckType = UserInfoUpdated
}

func (a *AckMessage) DataSend() {
	a.Clear()
	a.Type = DataAck
	a.DataAckType = Send
	a.IndexNumber = 10
	a.UserId = 1024
	a.Room = false
	a.RoomId = 0
}

func (a *AckMessage) DataRecieve() {
	a.Clear()
	a.Type = DataAck
	a.DataAckType = Receive
	a.IndexNumber = 1024
	a.UserId = 0
	a.Room = true
	a.RoomId = 0
}

func (a *AckMessage) DataRead() {
	a.Clear()
	a.Type = DataAck
	a.DataAckType = Read
	a.IndexNumber = 255
	a.UserId = 10
	a.Room = false
	a.RoomId = 0
}

func (a *AckMessage) AckMessageTestCoder(command uint8) []byte {
	switch command {
	case 1:
		a.ControlUserUpdate()
	case 2:
		a.DataSend()
	case 3:
		a.DataRecieve()
	case 4:
		a.DataRead()
	default:
		log.Println(`Unknown test for Data Message`)
		return nil
	}
	return Coder(a)
}
