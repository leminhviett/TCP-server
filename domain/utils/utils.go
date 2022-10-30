package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

type Message struct {
	ApplicationRoute string
	ApplicationData MessagePayload
}
type MessagePayload []byte


func WriteTo(conn net.Conn, message *Message) (int, error) {
	// 1. Convert data to byte
	messageB, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	var dataLength uint64 = uint64(len(messageB))


	// 2. Put data length & real data into place holder
	buf := make([]byte, 8+dataLength)
	// put data length to buffer
	binary.LittleEndian.PutUint64(buf, dataLength)
	//copy real data to buffer
	for idx := 8; idx < len(buf); idx ++ {
		buf[idx] = messageB[idx - 8]
	}
	
	// 3. Write data
	return conn.Write(buf)
}

func ReadFrom(conn net.Conn) (*Message, error){
	// 1. Get data length
	dataLengthB := make([]byte, 8)
	_, err := conn.Read(dataLengthB)
	if err != nil {
		panic(err)
	}
	dataLength := binary.LittleEndian.Uint64(dataLengthB)

	// 2. Get real data
	dataB := make([]byte, dataLength)
	_, err = conn.Read(dataB)
	if err != nil {
		panic(err)
	}

	fmt.Println(dataB)

	// 3. Return real data
	message := &Message{}
	err = json.Unmarshal(dataB, message)
	return message, err
}
