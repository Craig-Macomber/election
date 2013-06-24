package msg

import (
	"encoding/binary"
	"fmt"
	"io"
)

const Service = ":6038"

var Endian binary.ByteOrder = binary.LittleEndian

func ReadType(r io.Reader) (Type, error) {
	var messageType Type
	err := binary.Read(r, Endian, &messageType)
	return messageType, err
}

func ReadBlock(r io.Reader, lengthMax uint16) ([]byte, error) {
	var length uint16
	if err := binary.Read(r, Endian, &length); err != nil {
		return nil, err
	}
	if length > lengthMax {
		return nil, fmt.Errorf("Length %d to high (max=%f)", length, lengthMax)
	}
	buff := make([]byte, length)
	if err := binary.Read(r, Endian, &buff); err != nil {
		return buff, err
	}
	return buff, nil
}
func WriteBlock(w io.Writer, messageType Type, data []byte) error {
	if len(data) >= 1<<16 {
		return fmt.Errorf("Tried to send too long of message")
	}
	if err := binary.Write(w, Endian, messageType); err != nil {
		return err
	}
	length := uint16(len(data))
	if err := binary.Write(w, Endian, length); err != nil {
		return err
	}
	_, err := w.Write(data)
	return err

}
