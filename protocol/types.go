package protocol

import (
	"encoding/binary"
	"math/big"
)

// numeric
func numericDecoder(in []byte) (int, []byte, error) {
	if len(in) == 0 {
		return 0, in, ErrorMessageTooShort
	}
	if (in[0] >> 7) == 0 {
		return int(in[0]), in[1:], nil
	}

	length := in[0] & 0b0111_1111
	if len(in) < int(length+1) || length > 4 {
		return 0, in, ErrorLengthDecode
	}
	switch length {
	case 1:
		return int(in[1]), in[2:], nil
	case 2:
		return int(binary.BigEndian.Uint16(in[1:3])), in[3:], nil
	case 3:
		return int(binary.BigEndian.Uint32([]byte{0x00, in[1], in[2], in[3]})), in[4:], nil
	case 4:
		return int(binary.BigEndian.Uint64(in[1:5])), in[5:], nil
	}
	return 0, in, ErrorDecode
}

func numericCoder(num int) []byte {
	out := make([]byte, 1, 5)
	if num < 128 {
		out[0] = uint8(num)
	} else {
		num := big.NewInt(int64(num)).Bytes()
		out[0] = uint8(len(num)) | 0b1000_0000
		out = append(out, num...)
	}
	return out
}

// TitleNumeric
func titleNumericDecoder(in []byte) (int, []byte, error) {
	if len(in) < 2 {
		return 0, nil, ErrorMessageTooShort
	}
	if (in[0] >> 7) == 0 {
		return int(in[1]), in[2:], nil
	}
	length := in[1]
	if int(length)+2 > len(in) {
		return 0, in, ErrorLengthDecode
	}
	switch length {
	case 2:
		return int(binary.BigEndian.Uint16(in[2:4])), in[4:], nil
	case 3:
		return int(binary.BigEndian.Uint32([]byte{0x00, in[2], in[3], in[4]})), in[5:], nil
	case 4:
		return int(binary.BigEndian.Uint64(in[2:6])), in[6:], nil
	}
	return 0, in, ErrorLengthDecode
}

func titleNumericCoder(number int, title uint8) []byte {
	if number == 0 {
		return []byte{title, 0}
	}
	var length int
	num := big.NewInt(int64(number)).Bytes()
	if len(num) > 1 {
		title = 128 | title // Set 1 bit in value 1
		length = len(num) + 2
	} else {
		length = 2
	}
	out := make([]byte, 1, length)
	copy(out, []byte{title})
	if len(num) > 1 {
		out = append(out, uint8(len(num)))
	}
	out = append(out, num...)
	return out
}

// []byte

func dataDecoder(in []byte) ([]byte, []byte, error) {
	if len(in) < 2 {
		return nil, in, ErrorMessageTooShort
	}
	var length int
	if in[0]>>7 == 0 {
		length = int(in[0])
		in = in[1:]
	} else {
		length = int(binary.BigEndian.Uint16([]byte{(in[0] & 0b0111_1111), in[1]}))
		in = in[2:]
	}
	if length > len(in) {
		return nil, in, ErrorLengthDecode
	}
	return in[:length], in[length:], nil
}

func dataCoder(in []byte) []byte {
	if in == nil {
		return nil
	}
	length := len(in)
	var (
		lenField int
		lenVal   []byte
	)
	if length > 32767 {
		return nil
	}
	if length < 128 {
		lenField = 1
		lenVal = []byte{uint8(length)}
	} else {
		lenField = 2
		lenVal = doubleByteLenghtCoder(uint16(length))
	}
	out := make([]byte, lenField, len(in)+lenField)
	copy(out, lenVal)
	out = append(out, in...)
	return out
}

// Кодирование длин в 2 байта
func doubleByteLenghtCoder(length uint16) []byte {
	out := make([]byte, 2)
	binary.BigEndian.PutUint16(out, length)
	out[0] = out[0] | 0b1000_0000
	return out
}
