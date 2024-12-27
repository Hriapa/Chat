package protocol

import "encoding/binary"

type ControlMessageType uint8

const (
	UserNameRequestMeesageType ControlMessageType = iota + 1
	UserNameResponseMessageType
	ConnectMessageType
	DisconnectMessageType
	RegistrationMessageType
	UserUpdateMessageType
	UserInfoRequestMessageType
	UserInfoResponseMessageType
	UsersListMessageType
	MessagesRequestMessageType
	NumberOfMessagesMessageType
)

type ControlMessage struct {
	Type             ControlMessageType
	UserId           int
	RoomId           int
	UserName         string
	UserInfo         *UserInfo
	UsersList        UsersList
	MessagesRequest  MessagesRequest
	NumberOfMessages int
}

func NewControlMessage() *ControlMessage {
	return &ControlMessage{
		UserInfo:  &UserInfo{},
		UsersList: newUsersList(),
	}
}

func (c *ControlMessage) Clear() {
	c.Type = 0
	c.UserId = 0
	c.UserName = ""
	c.UserInfo.clear()
	c.UsersList.clear()
	c.MessagesRequest.clear()
	c.NumberOfMessages = 0
}

func (c *ControlMessage) Code() []byte {
	out := make([]byte, 2)
	copy(out, []byte{uint8(ControlMessageTitle), uint8(c.Type)})
	switch c.Type {
	case UserNameResponseMessageType:
		out = append(out, dataCoder([]byte(c.UserName))...)
	case ConnectMessageType, DisconnectMessageType:
		out = append(out, numericCoder(c.UserId)...)
	case RegistrationMessageType, UserUpdateMessageType:
		if c.UserName == "" {
			c.UserName = "Unknown"
		}
		out = append(out, numericCoder(c.UserId)...)
		out = append(out, dataCoder([]byte(c.UserName))...)
	case UserInfoResponseMessageType:
		if c.UserInfo != nil {
			out = append(out, c.UserInfo.codeUserInfo()...)
		}
	case UsersListMessageType:
		out = append(out, c.UsersList.codeUserList()...)
	case MessagesRequestMessageType:
		out = append(out, numericCoder(c.RoomId)...)
		out = append(out, c.MessagesRequest.codeMessagesRequest()...)
	case NumberOfMessagesMessageType:
		out = append(out, numericCoder(c.RoomId)...)
		out = append(out, numericCoder(c.NumberOfMessages)...)
	}
	return out
}

func (c *ControlMessage) Decode(in []byte) error {
	if len(in) < 2 {
		return ErrorMessageTooShort
	}
	if in[0] != byte(ControlMessageTitle) {
		return ErrorInMessageType
	}
	var (
		err  error
		name []byte
	)
	c.Type = ControlMessageType(in[1])
	in = in[2:]
	switch c.Type {
	case UserNameResponseMessageType:
		name, _, err = dataDecoder(in)
		c.UserName = string(name)
	case ConnectMessageType, DisconnectMessageType:
		c.UserId, _, err = numericDecoder(in)
	case RegistrationMessageType, UserUpdateMessageType:
		if len(in) < 3 {
			return ErrorLengthDecode
		}
		c.UserId, in, err = numericDecoder(in)
		if err != nil {
			return err
		}
		name, _, err = dataDecoder(in)
		c.UserName = string(name)
	case UserInfoResponseMessageType:
		if c.UserInfo == nil {
			c.UserInfo = &UserInfo{}
		}
		err = c.UserInfo.decodeUserInfo(in)
	case UsersListMessageType:
		if len(in) < 3 {
			return ErrorDecode
		}
		err = c.UsersList.decodeUserList(in)
	case MessagesRequestMessageType:
		c.RoomId, in, err = numericDecoder(in)
		if err != nil {
			return err
		}
		err = c.MessagesRequest.decodeMessagesRequest(in)
	case NumberOfMessagesMessageType:
		c.RoomId, in, err = numericDecoder(in)
		if err != nil {
			return err
		}
		c.NumberOfMessages, _, err = numericDecoder(in)
	}
	return err
}

// User Info

type UserInfo struct {
	Name       string
	Surname    string
	FamilyName string
	BirthDate  BirthDate
}

type BirthDate struct {
	Year  uint16
	Month uint8
	Day   uint8
}

func (u *UserInfo) clear() {
	u.Name = ""
	u.Surname = ""
	u.FamilyName = ""
	u.BirthDate.Day = 0
	u.BirthDate.Month = 0
	u.BirthDate.Year = 0
}

func (u *UserInfo) codeUserInfo() []byte {
	length := 11 + len(u.Name) + len(u.Surname) + len(u.FamilyName) // birtday = 1 day + 1 mounth + 2 year, 4 - titles, 3 length
	out := make([]byte, 0, length)
	if u.Name != "" {
		out = append(out, []byte{uint8(nameTitle)}...)
		out = append(out, dataCoder([]byte(u.Name))...)
	}
	if u.Surname != "" {
		out = append(out, []byte{uint8(surnameTitle)}...)
		out = append(out, dataCoder([]byte(u.Surname))...)
	}
	if u.FamilyName != "" {
		out = append(out, []byte{uint8(familyNameTitle)}...)
		out = append(out, dataCoder([]byte(u.FamilyName))...)
	}

	out = append(out, []byte{uint8(birthDateTitle), u.BirthDate.Day, u.BirthDate.Month}...)

	year := make([]byte, 2)
	binary.BigEndian.PutUint16(year, u.BirthDate.Year)
	out = append(out, year...)

	return out
}

func (u *UserInfo) decodeUserInfo(in []byte) error {
	if len(in) < 2 {
		return ErrorMessageTooShort
	}
	var (
		title Title
		err   error
		name  []byte
	)
	for {
		if len(in) < 2 {
			break
		}
		title = Title(in[0] & 0b0111_1111)
		switch title {
		case nameTitle:
			name, in, err = dataDecoder(in[1:])
			u.Name = string(name)
			if err != nil {
				return err
			}
		case surnameTitle:
			name, in, err = dataDecoder(in[1:])
			u.Surname = string(name)
			if err != nil {
				return err
			}
		case familyNameTitle:
			name, in, err = dataDecoder(in[1:])
			u.FamilyName = string(name)
			if err != nil {
				return err
			}
		case birthDateTitle:
			if len(in) < 5 {
				return ErrorLengthDecode
			}
			u.BirthDate.Day = in[1]
			u.BirthDate.Month = in[2]
			u.BirthDate.Year = binary.BigEndian.Uint16(in[3:5])
			in = in[5:]
		default:
			return ErrorTitleIncorrect
		}
	}
	return nil
}

// Users LIst

type UserListCommand uint8

// 0 - reserved
const (
	SetUserList UserListCommand = iota + 1
	AddToUserList
	RefreshUserList
)

type UsersList struct {
	UserListCommand  UserListCommand
	NumberOfElements int
	List             []User
}

type User struct {
	Id              int
	Name            string
	Online          bool
	NotReadMessages int
}

func newUsersList() UsersList {
	return UsersList{
		UserListCommand:  0,
		NumberOfElements: 0,
		List:             make([]User, 0, 100),
	}
}

func (u *UsersList) clear() {
	u.UserListCommand = 0
	u.NumberOfElements = 0
	u.List = u.List[:0]
}

func (u *UsersList) codeUserList() []byte {
	out := make([]byte, 1, 10000)
	copy(out, []byte{uint8(u.UserListCommand)})
	if u.UserListCommand == SetUserList {
		out = append(out, titleNumericCoder(u.NumberOfElements, uint8(numberListElementsTitle))...)
	}
	out = append(out, []byte{uint8(listElementsTitle)}...)
	for i := range u.List {
		out = append(out, numericCoder(u.List[i].Id)...)
		out = append(out, dataCoder([]byte(u.List[i].Name))...)
		out = append(out, 0x00)
		if u.List[i].Online {
			out[len(out)-1] = out[len(out)-1] | 0xf0
		}
		if u.List[i].NotReadMessages != 0 {
			out[len(out)-1] = out[len(out)-1] | 0x0f
			out = append(out, numericCoder(u.List[i].NotReadMessages)...)
		}
	}
	return out
}

func (u *UsersList) decodeUserList(in []byte) error {
	if len(in) < 2 {
		return ErrorMessageTooShort
	}
	var (
		err  error
		user User
		name []byte
	)
	if in[0]&0b1111_0000 != 0 { // проверка spare бит
		return ErrorDecode
	}
	u.UserListCommand = UserListCommand(in[0])
	if in[1]&0b0111_1111 == uint8(numberListElementsTitle) {
		u.NumberOfElements, in, err = titleNumericDecoder(in[1:])
		if err != nil {
			return err
		}
	}
	if len(in) == 0 {
		return nil
	}
	if in[0]&0b0111_1111 != uint8(listElementsTitle) {
		return ErrorTitleIncorrect
	}
	in = in[1:]
	for {
		if len(in) < 4 {
			break
		}
		user.Id, in, err = numericDecoder(in)
		if err != nil {
			return err
		}
		name, in, err = dataDecoder(in)
		user.Name = string(name)
		if err != nil {
			return err
		}
		user.Online = in[0]>>4 == 0x0f
		if in[0]&0b0000_1111 == 0x0f {
			user.NotReadMessages, in, err = numericDecoder(in[1:])
			if err != nil {
				return err
			}
		} else {
			user.NotReadMessages = 0
			in = in[1:]
		}
		u.List = append(u.List, user)
	}

	return nil
}

// Message Request

type MessagesRequest struct {
	LastMessageNumber int
}

func (m *MessagesRequest) clear() {
	m.LastMessageNumber = 0
}

func (m *MessagesRequest) codeMessagesRequest() []byte {
	return numericCoder(m.LastMessageNumber)
}

func (m *MessagesRequest) decodeMessagesRequest(in []byte) error {
	if len(in) == 0 {
		return ErrorMessageTooShort
	}
	var err error
	m.LastMessageNumber, _, err = numericDecoder(in)
	return err
}
