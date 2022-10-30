package utils

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUtils(t *testing.T){
	//https://stackoverflow.com/questions/69692663/how-to-unit-test-a-net-conn-function-that-modifies-the-message-sent
	serverConn, clientConn := net.Pipe()
	message := &Message{
		ApplicationRoute: "/hello",
		ApplicationData: []byte{1,2,3},
	}

	handleMessage := func() {
		_, err := WriteTo(serverConn, message)
		assert.NoError(t, err)
	}

	// Set deadline so test can detect if HandleMessage does not return
	clientConn.SetDeadline(time.Now().Add(time.Second))

	// Configure a go routine to act as the server
	go func() {
		handleMessage()
		serverConn.Close()
	}()

	getMessage, err := ReadFrom(clientConn)
	assert.NoError(t, err)
	assert.Equal(t, message, getMessage)
}

