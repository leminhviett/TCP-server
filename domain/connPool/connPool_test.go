package connpool

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/leminhviett/TCP-server/domain/customError"
	"github.com/stretchr/testify/assert"
)

func TestConnPool(t *testing.T) {
	ctx := context.Background()
	reset := setAndResetCreateNewConn()
	defer reset()

	var maxIdleConn int32 = 100
	var maxOpenConn int32 = 120

	pool := NewConnPool(ctx, maxIdleConn, maxOpenConn)
	waitgrp := sync.WaitGroup{}
	tempConnHolder := make([]*NetConn, 0)
	lock := sync.Mutex{}

	// 1. Test GetConn
	waitgrp.Add(int(maxOpenConn))
	for i := 0; i < int(maxOpenConn); i++ {
		go func() {
			con, err := pool.GetConn(ctx)
			lock.Lock()
			tempConnHolder = append(tempConnHolder, con)
			lock.Unlock()

			assert.NotNil(t, con)
			assert.NoError(t, err)
			waitgrp.Done()
		}()
	}
	waitgrp.Wait()
	assert.Equal(t, int32(120), pool.numberOpenConn)
	assert.Equal(t, 120, len(tempConnHolder))

	// 2. Test Release Conn
	waitgrp.Add(int(maxOpenConn))
	for i := 0; i < int(maxOpenConn); i++ {
		go func(idx int) {
			pool.PutConn(ctx, tempConnHolder[idx])
			waitgrp.Done()
		}(i)
	}
	waitgrp.Wait()
	assert.Equal(t, int(maxIdleConn), len(pool.freeConn))

	// 4. Test Queue
	for i := 0; i < 200; i++ {
		go func() {
			ctx, cancel := context.WithTimeout(ctx, 600*time.Microsecond)
			defer cancel()

			_, err := pool.GetConn(ctx) // expect put into waiting queue
			if err != nil {
				assert.ErrorIs(t, customError.ERR_GET_CONN_TIMEOUT, err)
			}
		}()
	}
	time.Sleep(600 * time.Microsecond)

	assert.Equal(t, int(maxOpenConn/2), len(pool.requestQueue))
	assert.Equal(t, 0, len(pool.freeConn))
}

func setAndResetCreateNewConn() func() {
	original := createNewConn

	conn, _ := net.Pipe()
	createNewConn = func() (*NetConn, error) {
		return &NetConn{
			Conn: conn,
		}, nil
	}

	return func() { createNewConn = original }

}
