package osuParser

import (
	"encoding/binary"
	"errors"
	"io"
	"time"
)

func readULEB128(r io.Reader) (uint64, error) {
	var result uint64
	var shift uint
	for {
		var byteVal byte
		if err := binary.Read(r, binary.LittleEndian, &byteVal); err != nil {
			return 0, err
		}
		result |= uint64(byteVal&0x7F) << shift
		if (byteVal & 0x80) == 0 {
			break
		}
		shift += 7
	}
	return result, nil
}

func readBoolean(r io.Reader) (bool, error) {
	var b byte
	if err := binary.Read(r, binary.LittleEndian, &b); err != nil {
		return false, err
	}
	return b != 0x00, nil
}

func readByte(r io.Reader) (byte, error) {
	var byte byte
	if err := binary.Read(r, binary.LittleEndian, &byte); err != nil {
		return 0x0b, err
	}

	return byte, nil
}

func readString(r io.Reader) (string, error) {
	var flag byte
	if err := binary.Read(r, binary.LittleEndian, &flag); err != nil {
		return "", err
	}

	if flag == 0x00 {
		return "", nil
	} else if flag == 0x0b {
		length, err := readULEB128(r)
		if err != nil {
			return "", err
		}
		data := make([]byte, length)
		if _, err := io.ReadFull(r, data); err != nil {
			return "", err
		}
		return string(data), nil
	}
	return "", errors.New("invalid string flag")
}

func readInt(r io.Reader) (int32, error) {
	var num int32
	if err := binary.Read(r, binary.LittleEndian, &num); err != nil {
		return 0, err
	}
	return num, nil
}

func readShort(r io.Reader) (uint16, error) {
	var num uint16
	if err := binary.Read(r, binary.LittleEndian, &num); err != nil {
		return 0, err
	}
	return num, nil
}

func readShortSigned(r io.Reader) (int16, error) {
	var num int16
	if err := binary.Read(r, binary.LittleEndian, &num); err != nil {
		return 0, err
	}
	return num, nil
}

func readLong(r io.Reader) (int64, error) {
	var num int64
	if err := binary.Read(r, binary.LittleEndian, &num); err != nil {
		return 0, err
	}
	return num, nil
}

func readSingle(r io.Reader) (float32, error) {
	var num float32
	if err := binary.Read(r, binary.LittleEndian, &num); err != nil {
		return 0, err
	}
	return num, nil
}

func readDouble(r io.Reader) (float64, error) {
	var num float64
	if err := binary.Read(r, binary.LittleEndian, &num); err != nil {
		return 0, err
	}
	return num, nil
}

func readDateTime(ticks int64) time.Time {
	const ticksPerSecond = 10000000
	const ticksOffset = 621355968000000000 // Ticks between 0001 and Unix epoch

	unixTicks := ticks - ticksOffset
	seconds := unixTicks / ticksPerSecond
	nanoseconds := (unixTicks % ticksPerSecond) * 100

	return time.Unix(seconds, nanoseconds).UTC()
}

func readIntDoublePairs(r io.Reader) (map[int]int64, error) {
	count, err := readInt(r)
	if err != nil {
		return nil, err
	}

	pairs := make(map[int]int64)
	for i := 0; i < int(count); i++ {
		var flag byte
		if err := binary.Read(r, binary.LittleEndian, &flag); err != nil {
			return nil, err
		}
		if flag != 0x08 {
			return nil, errors.New("invalid Int-Double pair flag")
		}

		intVal, err := readInt(r)
		if err != nil {
			return nil, err
		}

		var doubleFlag byte
		if err := binary.Read(r, binary.LittleEndian, &doubleFlag); err != nil {
			return nil, err
		}
		if doubleFlag != 0x0d {
			return nil, errors.New("invalid Double flag in Int-Double pair")
		}

		doubleVal, err := readDouble(r)
		if err != nil {
			return nil, err
		}

		pairs[int(intVal)] = int64(doubleVal)
	}
	return pairs, nil
}
