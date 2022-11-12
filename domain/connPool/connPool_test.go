package connpool

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConnPool(t *testing.T){
	ctx := context.Background()
	reset := setAndResetCreateNewConn()
	defer reset()

	pool := NewConnPool(100, 120)
	waitgrp := sync.WaitGroup{}
	tempConnHolder := make([]*NetConn, 0)
	lock := sync.Mutex{}

	// 1. Test GetConn
	waitgrp.Add(100)
	for i := 0; i < 100; i ++ {
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
	assert.Equal(t, int32(100), pool.openConn)
	assert.Equal(t, 120, len(tempConnHolder))

	// 2. Test GetConn
	waitgrp.Add(20)
	for i := 0; i < 20; i ++ {
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
	assert.Equal(t, int32(120), pool.openConn)
	assert.Equal(t, 120, len(tempConnHolder))


	// 3. Test Release Conn
	waitgrp.Add(120)
	for i := 0; i < 120; i ++ {
		go func(idx int) {
			pool.PutConn(ctx, tempConnHolder[idx])
			waitgrp.Done()
		}(i)
	}
	waitgrp.Wait()
	assert.Equal(t, 100, len(pool.freeConns))


	// 4. Test Queue
	for i := 0; i < 160; i ++ {
		go func() {
			con, err := pool.GetConn(ctx) // put into waiting queue
			assert.NotNil(t, con)
			assert.NoError(t, err)
		}()
	}
	time.Sleep(1*time.Second)

	conn, err := pool.GetConn(ctx)
	assert.Nil(t, conn)
	assert.ErrorIs(t, CONN_POOL_QUEUE_FULL, err)

	assert.Equal(t, 50, len(pool.requestQueue))
	assert.Equal(t, 0, len(pool.freeConns))
}

func setAndResetCreateNewConn() func(){
	original := createNewConn
	createNewConn = func() (*NetConn, error){
		return &NetConn{}, nil
	}

	return func() {createNewConn = original}

}