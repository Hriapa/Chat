package chat

import (
	"errors"
	"sort"
	"sync"
	"time"
	"webServer/protocol"
)

var (
	ErrorNoNillAssemblerBuffer = errors.New(`error: assembling buffer not nil`)
	ErrorNoMessagesInBuffer    = errors.New(`error: no messages in assembling`)
	ErrorAssembleMessage       = errors.New(`error: assembling message fail`)
	ErrorLostFirstPacket       = errors.New(`error: first packet losted`)
)

const fragmentLive = 30 // Время жизни фрагмента в секундах

type ErrorAssembling struct {
	Id     int
	Number int
}

type MessageAssembler struct {
	sync.Mutex
	Err      chan ErrorAssembling
	Messages map[int]*Message
}

type Message struct {
	timeStamp          time.Time
	FirstPacketInclude bool
	MessageNumber      int
	Fragments          []MessgeFragment
}

type MessgeFragment struct {
	SequenceNumber uint8
	Data           []byte
}

func NewMessageAssembler() *MessageAssembler {
	return &MessageAssembler{
		Err:      make(chan ErrorAssembling),
		Messages: make(map[int]*Message),
	}
}

func (m *MessageAssembler) Processor(message *protocol.DataMessage) (bool, []byte) {
	m.Lock()
	defer m.Unlock()

	var (
		id, previosNumber int
		out               []byte
		err               error
	)
	fragment := MessgeFragment{
		SequenceNumber: message.UserData.Fragmentation.Counter,
		Data:           message.UserData.Data,
	}
	if message.Room {
		id = message.RoomId
	} else {
		id = message.UserId
	}
	switch message.UserData.Fragmentation.FragmentType {
	case protocol.FirstFragment:
		previosNumber, err = m.startAssembling(id, message.IndexNumber, fragment)
		if err != nil {
			m.Err <- ErrorAssembling{id, previosNumber}
		}
		return false, nil
	case protocol.MiidleFragment:
		previosNumber, err = m.addAssembling(id, message.IndexNumber, fragment)
		if err != nil {
			m.Err <- ErrorAssembling{id, previosNumber}
		}
		return false, nil
	case protocol.LastFragment:
		out, err = m.completeAssembling(id, message.IndexNumber, message.UserData.Data)
		if err != nil {
			m.Err <- ErrorAssembling{id, previosNumber}
			return false, out
		}
		return true, out
	}
	return false, nil
}

func (m *MessageAssembler) startAssembling(id int, messageNumber int, fragment MessgeFragment) (int, error) {

	if val, ok := m.Messages[id]; ok {
		// Проверка на наличие в буфере не собранных сообщений. Вернёт ошибку по предыдущему сообщению
		if m.Messages[id].FirstPacketInclude || m.Messages[id].MessageNumber != messageNumber {
			prevMessageNumber := m.Messages[id].MessageNumber
			delete(m.Messages, id)
			m.createNewSession(id, messageNumber, true, fragment)
			return prevMessageNumber, ErrorNoNillAssemblerBuffer
		}
		val.timeStamp = time.Now()
		val.FirstPacketInclude = true
		val.Fragments = append(val.Fragments, fragment)
		return 0, nil
	}
	m.createNewSession(id, messageNumber, true, fragment)
	return 0, nil
}

func (m *MessageAssembler) createNewSession(id int, messageNumber int, firstPacket bool, fragment MessgeFragment) {
	data := &Message{
		timeStamp:          time.Now(),
		FirstPacketInclude: firstPacket,
		MessageNumber:      messageNumber,
		Fragments:          make([]MessgeFragment, 0, 3),
	}
	data.Fragments = append(data.Fragments, fragment)
	m.Messages[id] = data
}

func (m *MessageAssembler) addAssembling(id int, messageNumber int, fragment MessgeFragment) (int, error) {
	val, ok := m.Messages[id]
	if !ok {
		m.createNewSession(id, messageNumber, false, fragment)
		return 0, nil
	}
	if val.MessageNumber != messageNumber {
		prevMessageNumber := val.MessageNumber
		delete(m.Messages, id)
		m.createNewSession(id, messageNumber, false, fragment)
		return prevMessageNumber, ErrorNoNillAssemblerBuffer
	}
	val.timeStamp = time.Now()
	val.Fragments = append(val.Fragments, fragment)
	m.Messages[id] = val
	return 0, nil
}

func (m *MessageAssembler) completeAssembling(id int, messageNumber int, data []byte) ([]byte, error) {
	val, ok := m.Messages[id]
	if !ok {
		return nil, ErrorNoMessagesInBuffer
	}
	if !val.FirstPacketInclude {
		delete(m.Messages, id)
		return nil, ErrorLostFirstPacket
	}
	if val.MessageNumber != messageNumber {
		delete(m.Messages, id)
		return nil, ErrorNoNillAssemblerBuffer
	}
	sort.Slice(val.Fragments, func(i, j int) bool {
		return val.Fragments[i].SequenceNumber < val.Fragments[j].SequenceNumber
	})
	length := len(val.Fragments) + 1
	out := make([]byte, 0, length)
	for i := range val.Fragments {
		out = append(out, val.Fragments[i].Data...)
	}
	if data != nil {
		out = append(out, data...)
	}
	delete(m.Messages, id)
	return out, nil
}

func (m *MessageAssembler) assemblerCleane() {
	m.Lock()
	defer m.Unlock()

	var d time.Duration

	for key, val := range m.Messages {
		d = time.Since(val.timeStamp)
		if d.Seconds() > fragmentLive {
			m.Err <- ErrorAssembling{Id: key, Number: val.MessageNumber}
			delete(m.Messages, key)
		}
	}
}

func MessageFragmentation(in []byte, chunkSize int) []MessgeFragment {
	if len(in) < chunkSize {
		return []MessgeFragment{{SequenceNumber: 1, Data: in}}
	}
	numOfFragment := len(in) / chunkSize

	if len(in)%chunkSize != 0 {
		numOfFragment += 1
	}

	out := make([]MessgeFragment, numOfFragment)

	i := 1
	for {
		if (i-1)*chunkSize > len(in) {
			break
		}
		if (len(in) - (i-1)*chunkSize) >= chunkSize {
			out[i-1] = MessgeFragment{
				SequenceNumber: uint8(i),
				Data:           in[chunkSize*(i-1) : chunkSize*i],
			}
		} else {
			out[i-1] = MessgeFragment{
				SequenceNumber: uint8(i),
				Data:           in[chunkSize*(i-1):],
			}
		}
		i++
	}
	return out
}
