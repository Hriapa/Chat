package protocol

import (
	"bytes"
	"reflect"
	"testing"
)

func TestTextMessageCoder(t *testing.T) {
	for _, test := range []struct {
		name  string
		input DataMessage
		want  []byte
	}{
		{
			name: "Test_1",
			input: DataMessage{
				Type:        NewMessage,
				Format:      Text,
				IndexNumber: 10,
				UserId:      10,
				RoomId:      0,
				UserData: UserData{
					Fragmentation: Fragmentation{
						On:           false,
						FragmentType: 0,
						Counter:      0,
					},
					Data: []byte(`Hello`),
				},
			},
			want: []byte{0x01, 0x01, 0x01, 0x0a, 0x01, 0x0a, 0x0b, 0x05, 0x48, 0x65, 0x6c, 0x6c, 0x6f},
		}, {
			name: `Test_2`,
			input: DataMessage{
				Type:        OldMessage,
				Format:      Text,
				IndexNumber: 0,
				UserId:      10,
				Room:        true,
				RoomId:      1024,
				UserData: UserData{
					Fragmentation: Fragmentation{
						On:           false,
						FragmentType: 0,
						Counter:      0,
					},
					Data: []byte(`broadcast`),
				},
			},
			want: []byte{0x01, 0x02, 0x01, 0x00, 0x01, 0x0a, 0x82, 0x2, 0x04, 0x00, 0x0b, 0x09, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74},
		},
		{
			name: `Test_Lengt2Byte`,
			input: DataMessage{
				Type:        UpdateMessage,
				Format:      Text,
				IndexNumber: 1024,
				UserId:      10,
				RoomId:      0,
				UserData: UserData{
					Fragmentation: Fragmentation{
						On:           false,
						FragmentType: 0,
						Counter:      0,
					},
					Data: []byte(`Если позволять себе шутить, люди не воспринимают тебя всерьёз. И эти самые люди не понимают, что есть многое, чего нельзя выдержать, если не шутить.`),
				},
			},
			want: []byte{0x01, 0x03, 0x01, 0x82, 0x04, 0x00, 0x01, 0x0a, 0x0b, 0x81, 0x0b,
				0xd0, 0x95, 0xd1, 0x81, 0xd0, 0xbb, 0xd0, 0xb8, 0x20, 0xd0, 0xbf, 0xd0, 0xbe, 0xd0, 0xb7, 0xd0, 0xb2, 0xd0, 0xbe, 0xd0, 0xbb, 0xd1, 0x8f, 0xd1, 0x82, 0xd1, 0x8c, 0x20, 0xd1, 0x81, 0xd0, 0xb5,
				0xd0, 0xb1, 0xd0, 0xb5, 0x20, 0xd1, 0x88, 0xd1, 0x83, 0xd1, 0x82, 0xd0, 0xb8, 0xd1, 0x82, 0xd1, 0x8c, 0x2c, 0x20, 0xd0, 0xbb, 0xd1, 0x8e, 0xd0, 0xb4, 0xd0, 0xb8, 0x20, 0xd0, 0xbd, 0xd0, 0xb5,
				0x20, 0xd0, 0xb2, 0xd0, 0xbe, 0xd1, 0x81, 0xd0, 0xbf, 0xd1, 0x80, 0xd0, 0xb8, 0xd0, 0xbd, 0xd0, 0xb8, 0xd0, 0xbc, 0xd0, 0xb0, 0xd1, 0x8e, 0xd1, 0x82, 0x20, 0xd1, 0x82, 0xd0, 0xb5, 0xd0, 0xb1,
				0xd1, 0x8f, 0x20, 0xd0, 0xb2, 0xd1, 0x81, 0xd0, 0xb5, 0xd1, 0x80, 0xd1, 0x8c, 0xd1, 0x91, 0xd0, 0xb7, 0x2e, 0x20, 0xd0, 0x98, 0x20, 0xd1, 0x8d, 0xd1, 0x82, 0xd0, 0xb8, 0x20, 0xd1, 0x81, 0xd0, 0xb0,
				0xd0, 0xbc, 0xd1, 0x8b, 0xd0, 0xb5, 0x20, 0xd0, 0xbb, 0xd1, 0x8e, 0xd0, 0xb4, 0xd0, 0xb8, 0x20, 0xd0, 0xbd, 0xd0, 0xb5, 0x20, 0xd0, 0xbf, 0xd0, 0xbe, 0xd0, 0xbd, 0xd0, 0xb8, 0xd0, 0xbc, 0xd0, 0xb0,
				0xd1, 0x8e, 0xd1, 0x82, 0x2c, 0x20, 0xd1, 0x87, 0xd1, 0x82, 0xd0, 0xbe, 0x20, 0xd0, 0xb5, 0xd1, 0x81, 0xd1, 0x82, 0xd1, 0x8c, 0x20, 0xd0, 0xbc, 0xd0, 0xbd, 0xd0, 0xbe, 0xd0, 0xb3, 0xd0, 0xbe,
				0xd0, 0xb5, 0x2c, 0x20, 0xd1, 0x87, 0xd0, 0xb5, 0xd0, 0xb3, 0xd0, 0xbe, 0x20, 0xd0, 0xbd, 0xd0, 0xb5, 0xd0, 0xbb, 0xd1, 0x8c, 0xd0, 0xb7, 0xd1, 0x8f, 0x20, 0xd0, 0xb2, 0xd1, 0x8b, 0xd0, 0xb4,
				0xd0, 0xb5, 0xd1, 0x80, 0xd0, 0xb6, 0xd0, 0xb0, 0xd1, 0x82, 0xd1, 0x8c, 0x2c, 0x20, 0xd0, 0xb5, 0xd1, 0x81, 0xd0, 0xbb, 0xd0, 0xb8, 0x20, 0xd0, 0xbd, 0xd0, 0xb5, 0x20, 0xd1, 0x88, 0xd1, 0x83,
				0xd1, 0x82, 0xd0, 0xb8, 0xd1, 0x82, 0xd1, 0x8c, 0x2e},
		},
		{
			name: `Test_Fragmentation_firstSegment`,
			input: DataMessage{
				Type:        NewMessage,
				Format:      Text,
				IndexNumber: 250,
				UserId:      10,
				RoomId:      0,
				UserData: UserData{
					Fragmentation: Fragmentation{
						On:           true,
						FragmentType: FirstFragment,
						Counter:      0,
					},
					Data: []byte(`first segment`),
				},
			},
			want: []byte{0x01, 0x01, 0x01, 0x81, 0xfa, 0x01, 0x0a, 0x8b, 0x80, 0x0d, 0x66, 0x69, 0x72, 0x73, 0x74, 0x20, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74},
		},
		{
			name: `Test_Fragmentation_lastSegment`,
			input: DataMessage{
				Type:        NewMessage,
				Format:      Text,
				IndexNumber: 10,
				UserId:      10,
				RoomId:      0,
				UserData: UserData{
					Fragmentation: Fragmentation{
						On:           true,
						FragmentType: LastFragment,
						Counter:      10,
					},
					Data: []byte(`last segment`),
				},
			},
			want: []byte{0x01, 0x01, 0x01, 0x0a, 0x01, 0x0a, 0x8b, 0x4a, 0x0c, 0x6c, 0x61, 0x73, 0x74, 0x20, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74},
		},
	} {
		d := &test.input
		result := Coder(d)
		if !bytes.Equal(result, test.want) {
			t.Logf("%s result is not expected \n want %v, \n got  %v", test.name, test.want, result)
			t.Fail()
		}
	}
}

func TestTextMessageDecoder(t *testing.T) {
	for _, test := range []struct {
		name  string
		input []byte
		want  DataMessage
	}{
		{
			name:  `Test1`,
			input: []byte{0x01, 0x01, 0x01, 0x0a, 0x01, 0x0a, 0x0b, 0x05, 0x48, 0x65, 0x6c, 0x6c, 0x6f},
			want: DataMessage{
				Type:        NewMessage,
				Format:      Text,
				IndexNumber: 10,
				UserId:      10,
				Room:        false,
				RoomId:      0,
				UserData: UserData{
					Fragmentation: Fragmentation{
						On:           false,
						FragmentType: 0,
						Counter:      0,
					},
					Data: []byte(`Hello`),
				},
			},
		},
		{
			name:  `Test2`,
			input: []byte{0x01, 0x02, 0x01, 0x00, 0x01, 0x0a, 0x82, 0x02, 0x04, 0x00, 0x0b, 0x09, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74},
			want: DataMessage{
				Type:        OldMessage,
				Format:      Text,
				IndexNumber: 0,
				UserId:      10,
				Room:        true,
				RoomId:      1024,
				UserData: UserData{
					Fragmentation: Fragmentation{
						On:           false,
						FragmentType: 0,
						Counter:      0,
					},
					Data: []byte(`broadcast`),
				},
			},
		},
		{
			name: `Test_Length2Byte`,
			input: []byte{0x01, 0x03, 0x01, 0x82, 0x04, 0x00, 0x01, 0x0a, 0x0b, 0x81, 0x0b,
				0xd0, 0x95, 0xd1, 0x81, 0xd0, 0xbb, 0xd0, 0xb8, 0x20, 0xd0, 0xbf, 0xd0, 0xbe, 0xd0, 0xb7, 0xd0, 0xb2, 0xd0, 0xbe, 0xd0, 0xbb, 0xd1, 0x8f, 0xd1, 0x82, 0xd1, 0x8c, 0x20, 0xd1, 0x81, 0xd0, 0xb5,
				0xd0, 0xb1, 0xd0, 0xb5, 0x20, 0xd1, 0x88, 0xd1, 0x83, 0xd1, 0x82, 0xd0, 0xb8, 0xd1, 0x82, 0xd1, 0x8c, 0x2c, 0x20, 0xd0, 0xbb, 0xd1, 0x8e, 0xd0, 0xb4, 0xd0, 0xb8, 0x20, 0xd0, 0xbd, 0xd0, 0xb5,
				0x20, 0xd0, 0xb2, 0xd0, 0xbe, 0xd1, 0x81, 0xd0, 0xbf, 0xd1, 0x80, 0xd0, 0xb8, 0xd0, 0xbd, 0xd0, 0xb8, 0xd0, 0xbc, 0xd0, 0xb0, 0xd1, 0x8e, 0xd1, 0x82, 0x20, 0xd1, 0x82, 0xd0, 0xb5, 0xd0, 0xb1,
				0xd1, 0x8f, 0x20, 0xd0, 0xb2, 0xd1, 0x81, 0xd0, 0xb5, 0xd1, 0x80, 0xd1, 0x8c, 0xd1, 0x91, 0xd0, 0xb7, 0x2e, 0x20, 0xd0, 0x98, 0x20, 0xd1, 0x8d, 0xd1, 0x82, 0xd0, 0xb8, 0x20, 0xd1, 0x81, 0xd0, 0xb0,
				0xd0, 0xbc, 0xd1, 0x8b, 0xd0, 0xb5, 0x20, 0xd0, 0xbb, 0xd1, 0x8e, 0xd0, 0xb4, 0xd0, 0xb8, 0x20, 0xd0, 0xbd, 0xd0, 0xb5, 0x20, 0xd0, 0xbf, 0xd0, 0xbe, 0xd0, 0xbd, 0xd0, 0xb8, 0xd0, 0xbc, 0xd0, 0xb0,
				0xd1, 0x8e, 0xd1, 0x82, 0x2c, 0x20, 0xd1, 0x87, 0xd1, 0x82, 0xd0, 0xbe, 0x20, 0xd0, 0xb5, 0xd1, 0x81, 0xd1, 0x82, 0xd1, 0x8c, 0x20, 0xd0, 0xbc, 0xd0, 0xbd, 0xd0, 0xbe, 0xd0, 0xb3, 0xd0, 0xbe,
				0xd0, 0xb5, 0x2c, 0x20, 0xd1, 0x87, 0xd0, 0xb5, 0xd0, 0xb3, 0xd0, 0xbe, 0x20, 0xd0, 0xbd, 0xd0, 0xb5, 0xd0, 0xbb, 0xd1, 0x8c, 0xd0, 0xb7, 0xd1, 0x8f, 0x20, 0xd0, 0xb2, 0xd1, 0x8b, 0xd0, 0xb4,
				0xd0, 0xb5, 0xd1, 0x80, 0xd0, 0xb6, 0xd0, 0xb0, 0xd1, 0x82, 0xd1, 0x8c, 0x2c, 0x20, 0xd0, 0xb5, 0xd1, 0x81, 0xd0, 0xbb, 0xd0, 0xb8, 0x20, 0xd0, 0xbd, 0xd0, 0xb5, 0x20, 0xd1, 0x88, 0xd1, 0x83,
				0xd1, 0x82, 0xd0, 0xb8, 0xd1, 0x82, 0xd1, 0x8c, 0x2e},
			want: DataMessage{
				Type:        UpdateMessage,
				Format:      Text,
				IndexNumber: 1024,
				UserId:      10,
				RoomId:      0,
				UserData: UserData{
					Fragmentation: Fragmentation{
						On:           false,
						FragmentType: 0,
						Counter:      0,
					},
					Data: []byte(`Если позволять себе шутить, люди не воспринимают тебя всерьёз. И эти самые люди не понимают, что есть многое, чего нельзя выдержать, если не шутить.`),
				},
			},
		},
		{
			name:  `Test_Fragmentation_firstSegment`,
			input: []byte{0x01, 0x01, 0x01, 0x81, 0xfa, 0x01, 0x0a, 0x8b, 0x80, 0x0d, 0x66, 0x69, 0x72, 0x73, 0x74, 0x20, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74},
			want: DataMessage{
				Type:        NewMessage,
				Format:      Text,
				IndexNumber: 250,
				UserId:      10,
				RoomId:      0,
				UserData: UserData{
					Fragmentation: Fragmentation{
						On:           true,
						FragmentType: FirstFragment,
						Counter:      0,
					},
					Data: []byte(`first segment`),
				},
			},
		},
		{
			name:  `Test_Fragmentation_lastSegment`,
			input: []byte{0x01, 0x01, 0x01, 0x0a, 0x01, 0x0a, 0x8b, 0x4a, 0x0c, 0x6c, 0x61, 0x73, 0x74, 0x20, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74},
			want: DataMessage{
				Type:        NewMessage,
				Format:      Text,
				IndexNumber: 10,
				UserId:      10,
				RoomId:      0,
				UserData: UserData{
					Fragmentation: Fragmentation{
						On:           true,
						FragmentType: LastFragment,
						Counter:      10,
					},
					Data: []byte(`last segment`),
				},
			},
		},
	} {
		d := DataMessage{}
		err := Decoder(&d, test.input)
		if err != nil {
			t.Errorf(`error in data message decoder \n error: %v`, err)
		}
		if !reflect.DeepEqual(d, test.want) {
			t.Logf("%s result is not expected \n want %v, \n got  %v", test.name, test.want, d)
			t.Fail()
		}
	}
}

func TestControlMessageCoder(t *testing.T) {
	for _, test := range []struct {
		name  string
		input ControlMessage
		want  []byte
	}{
		{
			name: `User_Name_Request`,
			input: ControlMessage{
				Type: UserNameRequestMeesageType,
			},
			want: []byte{0x00, 0x01},
		},
		{
			name: `User_Name_Response`,
			input: ControlMessage{
				Type:     UserNameResponseMessageType,
				UserName: `Ivan`,
			},
			want: []byte{0x00, 0x02, 0x04, 0x49, 0x76, 0x61, 0x6e},
		},
		{
			name: `Test_Connect`,
			input: ControlMessage{
				Type:   ConnectMessageType,
				UserId: 10,
			},
			want: []byte{0x00, 0x03, 0x0a},
		},
		{
			name: `Test_Disconnect`,
			input: ControlMessage{
				Type:   DisconnectMessageType,
				UserId: 1024,
			},
			want: []byte{0x00, 0x04, 0x82, 0x04, 0x00},
		},
		{
			name: `Test_Registration`,
			input: ControlMessage{
				Type:     RegistrationMessageType,
				UserId:   66000,
				UserName: "ivan",
			},
			want: []byte{0x00, 0x05, 0x83, 0x01, 0x01, 0xd0, 0x04, 0x69, 0x76, 0x61, 0x6e},
		},
		{
			name: `Test_UserUpdate`,
			input: ControlMessage{
				Type:     UserUpdateMessageType,
				UserId:   1024,
				UserName: `Ivan`,
			},
			want: []byte{0x00, 0x06, 0x82, 0x04, 0x00, 0x04, 0x49, 0x76, 0x61, 0x6e},
		},
		{
			name: "Test_UserInfo_Request",
			input: ControlMessage{
				Type: UserInfoRequestMessageType,
			},
			want: []byte{0x00, 0x07},
		},
		{
			name: "Test_UserInfo_Response",
			input: ControlMessage{
				Type: UserInfoResponseMessageType,
				UserInfo: &UserInfo{
					Name:       `Ivan`,
					Surname:    `Ivanovich`,
					FamilyName: `Ivanov`,
					BirthDate: BirthDate{
						Year:  2006,
						Month: 02,
						Day:   07,
					},
				},
			},
			want: []byte{0x00, 0x08, 0x03, 0x04, 0x49, 0x76, 0x61, 0x6e, 0x04, 0x09, 0x49, 0x76, 0x61, 0x6e, 0x6f, 0x76, 0x69, 0x63, 0x68,
				0x05, 0x06, 0x49, 0x76, 0x61, 0x6e, 0x6f, 0x76, 0x06, 0x07, 0x02, 0x07, 0xd6},
		},
		{
			name: "Test_UserInfo_Response_2",
			input: ControlMessage{
				Type: UserInfoResponseMessageType,
				UserInfo: &UserInfo{
					Name:       "",
					Surname:    "",
					FamilyName: "",
					BirthDate: BirthDate{
						Year:  2006,
						Month: 02,
						Day:   07,
					},
				},
			},
			want: []byte{0x00, 0x08, 0x06, 0x07, 0x02, 0x07, 0xd6},
		},
		{
			name: "Test_UserInfo_Update",
			input: ControlMessage{
				Type: UserInfoUpdateMessageType,
				UserInfo: &UserInfo{
					Name:       `Ivan`,
					Surname:    `Ivanovich`,
					FamilyName: ``,
					BirthDate: BirthDate{
						Year:  0,
						Month: 0,
						Day:   0,
					},
				},
			},
			want: []byte{0x00, 0x09, 0x03, 0x04, 0x49, 0x76, 0x61, 0x6e, 0x04, 0x09, 0x49, 0x76, 0x61, 0x6e, 0x6f, 0x76, 0x69, 0x63, 0x68},
		},
		{
			name: "Test_UserList",
			input: ControlMessage{
				Type: UsersListMessageType,
				UsersList: UsersList{
					UserListCommand:  SetUserList,
					NumberOfElements: 3,
					List: []User{
						{Id: 1, Name: "Petr", Online: false, NotReadMessages: 1},
						{Id: 1024, Name: "Ivan", Online: true, NotReadMessages: 4},
						{Id: 10, Name: "Olga", Online: true, NotReadMessages: 0},
					},
				},
			},
			want: []byte{0x00, 0x0a, 0x01, 0x09, 0x03, 0x08, 0x01, 0x04, 0x50, 0x65, 0x74, 0x72, 0x0f, 0x01, 0x82, 0x04, 0x00, 0x04,
				0x49, 0x76, 0x61, 0x6e, 0xff, 0x04, 0x0a, 0x04, 0x4f, 0x6c, 0x67, 0x61, 0xf0},
		},
		{
			name: "Test_MessagesRequest",
			input: ControlMessage{
				Type:   MessagesRequestMessageType,
				RoomId: 1024,
				MessagesRequest: MessagesRequest{
					LastMessageNumber: 256,
				},
			},
			want: []byte{0x00, 0x0b, 0x82, 0x04, 0x00, 0x82, 0x01, 0x00},
		},
		{
			name: "Test_NumberOfMessages",
			input: ControlMessage{
				Type:             NumberOfMessagesMessageType,
				RoomId:           10,
				NumberOfMessages: 256,
			},
			want: []byte{0x00, 0x0c, 0x0a, 0x82, 0x01, 0x00},
		},
	} {
		c := &test.input
		res := Coder(c)
		if !bytes.Equal(res, test.want) {
			t.Logf("%s result is not expected \n want %v, \n got  %v", test.name, test.want, res)
			t.Fail()
		}
	}
}

func TestControlMessageDecoder(t *testing.T) {
	for _, test := range []struct {
		name  string
		input []byte
		want  ControlMessage
	}{
		{
			name:  `User_Name_Request`,
			input: []byte{0x00, 0x01},
			want: ControlMessage{
				Type: UserNameRequestMeesageType,
			},
		},
		{
			name:  `User_Name_Response`,
			input: []byte{0x00, 0x02, 0x04, 0x49, 0x76, 0x61, 0x6e},
			want: ControlMessage{
				Type:     UserNameResponseMessageType,
				UserName: `Ivan`,
			},
		},
		{
			name:  `Test_Connect`,
			input: []byte{0x00, 0x03, 0x0a},
			want: ControlMessage{
				Type:   ConnectMessageType,
				UserId: 10,
			},
		},
		{
			name:  `Test_Disconnect`,
			input: []byte{0x00, 0x04, 0x82, 0x04, 0x00},
			want: ControlMessage{
				Type:   DisconnectMessageType,
				UserId: 1024,
			},
		},
		{
			name:  `Test_Registration`,
			input: []byte{0x00, 0x05, 0x83, 0x01, 0x01, 0xd0, 0x04, 0x69, 0x76, 0x61, 0x6e},
			want: ControlMessage{
				Type:     RegistrationMessageType,
				UserId:   66000,
				UserName: "ivan",
			},
		},
		{
			name:  `Test_UserUpdate`,
			input: []byte{0x00, 0x06, 0x82, 0x04, 0x00, 0x04, 0x49, 0x76, 0x61, 0x6e},
			want: ControlMessage{
				Type:     UserUpdateMessageType,
				UserId:   1024,
				UserName: "Ivan",
			},
		},
		{
			name:  `Test_User_Info_Reqest`,
			input: []byte{0x00, 0x07},
			want: ControlMessage{
				Type: UserInfoRequestMessageType,
			},
		},
		{
			name: `Test_UserInfo`,
			input: []byte{0x00, 0x08, 0x03, 0x04, 0x49, 0x76, 0x61, 0x6e, 0x04, 0x09, 0x49, 0x76, 0x61, 0x6e, 0x6f, 0x76, 0x69, 0x63, 0x68,
				0x05, 0x06, 0x49, 0x76, 0x61, 0x6e, 0x6f, 0x76, 0x06, 0x07, 0x02, 0x07, 0xd6},
			want: ControlMessage{
				Type: UserInfoResponseMessageType,
				UserInfo: &UserInfo{
					Name:       `Ivan`,
					Surname:    `Ivanovich`,
					FamilyName: `Ivanov`,
					BirthDate: BirthDate{
						Year:  2006,
						Month: 02,
						Day:   07,
					},
				},
			},
		},
		{
			name:  "Test_UserInfo_2",
			input: []byte{0x00, 0x08, 0x06, 0x07, 0x02, 0x07, 0xd6},
			want: ControlMessage{
				Type: UserInfoResponseMessageType,
				UserInfo: &UserInfo{
					Name:       "",
					Surname:    "",
					FamilyName: "",
					BirthDate: BirthDate{
						Year:  2006,
						Month: 02,
						Day:   07,
					},
				},
			},
		},
		{
			name:  "Test_UserInfo_Update",
			input: []byte{0x00, 0x09, 0x03, 0x04, 0x49, 0x76, 0x61, 0x6e, 0x04, 0x09, 0x49, 0x76, 0x61, 0x6e, 0x6f, 0x76, 0x69, 0x63, 0x68},
			want: ControlMessage{
				Type: UserInfoUpdateMessageType,
				UserInfo: &UserInfo{
					Name:       `Ivan`,
					Surname:    `Ivanovich`,
					FamilyName: ``,
					BirthDate: BirthDate{
						Year:  0,
						Month: 0,
						Day:   0,
					},
				},
			},
		},
		{
			name: "Test_UsersList",
			input: []byte{0x00, 0x0a, 0x01, 0x09, 0x03, 0x08, 0x01, 0x04, 0x50, 0x65, 0x74, 0x72, 0x0f, 0x01, 0x82, 0x04, 0x00, 0x04,
				0x49, 0x76, 0x61, 0x6e, 0xff, 0x04, 0x0a, 0x04, 0x4f, 0x6c, 0x67, 0x61, 0xf0},
			want: ControlMessage{
				Type: UsersListMessageType,
				UsersList: UsersList{
					UserListCommand:  SetUserList,
					NumberOfElements: 3,
					List: []User{
						{Id: 1, Name: "Petr", Online: false, NotReadMessages: 1},
						{Id: 1024, Name: "Ivan", Online: true, NotReadMessages: 4},
						{Id: 10, Name: "Olga", Online: true, NotReadMessages: 0},
					},
				},
			},
		},
		{
			name:  "Test_MessagesRequest",
			input: []byte{0x00, 0x0b, 0x82, 0x04, 0x00, 0x82, 0x01, 0x00},
			want: ControlMessage{
				Type:   MessagesRequestMessageType,
				RoomId: 1024,
				MessagesRequest: MessagesRequest{
					LastMessageNumber: 256,
				},
			},
		},
		{
			name:  "Test_NumberOfMessages",
			input: []byte{0x00, 0x0c, 0x0a, 0x82, 0x01, 0x00},
			want: ControlMessage{
				Type:             NumberOfMessagesMessageType,
				RoomId:           10,
				NumberOfMessages: 256,
			},
		},
	} {
		c := ControlMessage{}
		err := Decoder(&c, test.input)
		if err != nil {
			t.Errorf(`error in control message decoder, error: %v`, err)
		}
		if !reflect.DeepEqual(c, test.want) {
			t.Logf("%s result is not expected \n want %v, \n got %v", test.name, test.want, c)
			t.Fail()
		}
	}
}

func TestAckMessageCoder(t *testing.T) {
	for _, test := range []struct {
		name  string
		input AckMessage
		want  []byte
	}{
		{
			name: `Test ControlAckUserUpdate`,
			input: AckMessage{
				Type:           ControlAck,
				ControlAckType: UserInfoUpdated,
			},
			want: []byte{0x02, 0x01, 0x01},
		},
		{
			name: `Test DataAckSend`,
			input: AckMessage{
				Type:        DataAck,
				DataAckType: Send,
				IndexNumber: 10,
				UserId:      1024,
				Room:        false,
				RoomId:      0,
			},
			want: []byte{0x02, 0x02, 0x01, 0x0a, 0x81, 0x02, 0x04, 0x00},
		},
		{
			name: `Test DataAckReceived`,
			input: AckMessage{
				Type:        DataAck,
				DataAckType: Receive,
				IndexNumber: 1024,
				UserId:      0,
				Room:        true,
				RoomId:      0,
			},
			want: []byte{0x02, 0x02, 0x02, 0x82, 0x04, 0x00, 0x02, 0x00},
		},
		{
			name: `Test DataAckRead`,
			input: AckMessage{
				Type:        DataAck,
				DataAckType: Read,
				IndexNumber: 255,
				UserId:      10,
				Room:        false,
				RoomId:      0,
			},
			want: []byte{0x02, 0x02, 0x03, 0x81, 0xff, 0x01, 0x0a},
		},
	} {
		ack := &test.input
		res := Coder(ack)
		if !reflect.DeepEqual(res, test.want) {
			t.Logf("%s result is not expected \n want %v, \n got %v", test.name, test.want, res)
			t.Fail()
		}
	}
}

func TestAckMessageDecoder(t *testing.T) {
	for _, test := range []struct {
		name  string
		input []byte
		want  AckMessage
	}{
		{
			name:  `Test ControlAckUserUpdate`,
			input: []byte{0x02, 0x01, 0x01},
			want: AckMessage{
				Type:           ControlAck,
				ControlAckType: UserInfoUpdated,
			},
		},
		{
			name:  `Test AckSend`,
			input: []byte{0x02, 0x02, 0x01, 0x0a, 0x81, 0x02, 0x04, 0x00},
			want: AckMessage{
				Type:        DataAck,
				DataAckType: Send,
				IndexNumber: 10,
				UserId:      1024,
				Room:        false,
				RoomId:      0,
			},
		},
		{
			name:  `Test AckReceived`,
			input: []byte{0x02, 0x02, 0x02, 0x82, 0x04, 0x00, 0x02, 0x00},
			want: AckMessage{
				Type:        DataAck,
				DataAckType: Receive,
				IndexNumber: 1024,
				UserId:      0,
				Room:        true,
				RoomId:      0,
			},
		},
		{
			name:  `Test AckRead`,
			input: []byte{0x02, 0x02, 0x03, 0x81, 0xff, 0x01, 0x0a},
			want: AckMessage{
				Type:        DataAck,
				DataAckType: Read,
				IndexNumber: 255,
				UserId:      10,
				Room:        false,
				RoomId:      0,
			},
		},
	} {
		ack := AckMessage{}
		err := Decoder(&ack, test.input)
		if err != nil {
			t.Errorf(`error in acknowledge message decoder, error: %v`, err)
		}
		if !reflect.DeepEqual(ack, test.want) {
			t.Logf("%s result is not expected \n want %v, \n got %v", test.name, test.want, ack)
			t.Fail()
		}
	}
}

func TestErrorMessageCoder(t *testing.T) {
	for _, test := range []struct {
		name  string
		input ErrorMessage
		want  []byte
	}{
		{
			name: `Test_ControlError`,
			input: ErrorMessage{
				Type:               ControlMessageError,
				ControlMessageType: UsersListMessageType,
			},
			want: []byte{0x03, 0x01, 0x08},
		},
		{
			name: `Test_DataError_1`,
			input: ErrorMessage{
				Type:        DataMessageError,
				IndexNumber: 10,
				UserId:      10,
				RoomId:      0,
			},
			want: []byte{0x03, 0x02, 0x0a, 0x01, 0x0a},
		},
		{
			name: `Test_DataError_2`,
			input: ErrorMessage{
				Type:        DataMessageError,
				IndexNumber: 1024,
				UserId:      0,
				RoomId:      1024,
			},
			want: []byte{0x03, 0x02, 0x82, 0x04, 0x00, 0x82, 0x02, 0x04, 0x00},
		},
	} {
		e := &test.input
		res := e.Code()
		if !reflect.DeepEqual(test.want, res) {
			t.Logf("%s result is not expected \n want %v, \n got %v", test.name, test.want, res)
			t.Fail()
		}
	}
}

func TestErrorMessageDecoder(t *testing.T) {
	for _, test := range []struct {
		name  string
		input []byte
		want  ErrorMessage
	}{
		{
			name:  `Test_ControlError`,
			input: []byte{0x03, 0x01, 0x08},
			want: ErrorMessage{
				Type:               ControlMessageError,
				ControlMessageType: UsersListMessageType,
			},
		},
		{
			name:  `Test_DataError_1`,
			input: []byte{0x03, 0x02, 0x0a, 0x01, 0x0a},
			want: ErrorMessage{
				Type:        DataMessageError,
				IndexNumber: 10,
				UserId:      10,
				RoomId:      0,
			},
		},
		{
			name:  `Test_DataError_2`,
			input: []byte{0x03, 0x02, 0x82, 0x04, 0x00, 0x82, 0x02, 0x04, 0x00},
			want: ErrorMessage{
				Type:        DataMessageError,
				IndexNumber: 1024,
				UserId:      0,
				RoomId:      1024,
			},
		},
	} {
		e := ErrorMessage{}
		err := e.Decode(test.input)
		if err != nil {
			t.Errorf(`error in error message decoder, error: %v`, err)
		}
		if !reflect.DeepEqual(e, test.want) {
			t.Logf("%s result is not expected \n want %v, \n got %v", test.name, test.want, e)
			t.Fail()
		}
	}
}
