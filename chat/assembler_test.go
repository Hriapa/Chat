package chat

import (
	"bytes"
	"reflect"
	"testing"
	"webServer/protocol"
)

func TestAssembler(t *testing.T) {
	assemblingData := make([]*protocol.DataMessage, 3)

	for i := range assemblingData {
		assemblingData[i] = &protocol.DataMessage{
			Type:        protocol.NewMessage,
			Format:      protocol.Text,
			IndexNumber: 10,
			UserId:      10,
			Room:        false,
			RoomId:      0,
		}
		switch i {
		case 0:
			assemblingData[i].UserData = protocol.UserData{
				Fragmentation: protocol.Fragmentation{
					On:           true,
					FragmentType: protocol.FirstFragment,
					Counter:      1,
				},
				Data: []byte(`1st fragment `),
			}
		case 1:
			assemblingData[i].UserData = protocol.UserData{
				Fragmentation: protocol.Fragmentation{
					On:           true,
					FragmentType: protocol.MiidleFragment,
					Counter:      2,
				},
				Data: []byte(`2nd fragment `),
			}
		case 2:
			assemblingData[i].UserData = protocol.UserData{
				Fragmentation: protocol.Fragmentation{
					On:           true,
					FragmentType: protocol.LastFragment,
					Counter:      3,
				},
				Data: []byte(`3rd fragment`),
			}
		}
	}

	var (
		result []byte
		ok     bool
	)

	assembler := NewMessageAssembler()

	for i := range assemblingData {
		ok, result = assembler.Processor(assemblingData[i])
	}

	go func() {
		message := <-assembler.Err
		t.Errorf("error in assembling store /n count number: %v", message.Number)
	}()

	if !ok {
		t.Errorf(`error assembling packet not complete`)
	}

	if !bytes.Equal(result, []byte(`1st fragment 2nd fragment 3rd fragment`)) {
		t.Logf("assemling result not correct \n result: %s", result)
		t.Fail()
	}

}

func TestMessageFragmentation(t *testing.T) {
	input := []byte(`1st fragment 2nd fragment 3rd fragment`)
	want := make([]MessgeFragment, 3)
	chunkSize := 13

	want[0] = MessgeFragment{
		SequenceNumber: 1,
		Data:           []byte(`1st fragment `),
	}

	want[1] = MessgeFragment{
		SequenceNumber: 2,
		Data:           []byte(`2nd fragment `),
	}

	want[2] = MessgeFragment{
		SequenceNumber: 3,
		Data:           []byte(`3rd fragment`),
	}

	result := MessageFragmentation(input, chunkSize)

	if !reflect.DeepEqual(result, want) {
		t.Errorf("fragmentation result incorrect \n want: %v \n result: %v", want, result)
		t.Fail()
	}
}
