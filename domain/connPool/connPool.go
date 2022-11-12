package connpool

import (
	"context"
	"errors"
	"net"
	"sync"
)

var (
	ERR_GET_CONN_TIMEOUT = errors.New("Error get conn timeout")
	ERR_FREE_CONN_CHAN_CLOSED = errors.New("Free conn channel error")
	CONN_POOL_QUEUE_FULL = errors.New("Conn q is full")
)

type NetConn struct {
	net.Conn
}

type ConnPool interface{
	GetConn(ctx context.Context) (*NetConn, error)
	PutConn(ctx context.Context, conn *NetConn)
}

type ConnPoolImpl struct {
	freeConns chan *NetConn
	requestQueue chan (chan *NetConn)
	maxOpenConn int32
	openConn int32
	lock sync.Mutex
}

func NewConnPool(maxIdleConn, maxOpenConn int32) *ConnPoolImpl {
	pool := &ConnPoolImpl{
		freeConns: make(chan *NetConn, maxIdleConn),
		requestQueue: make(chan (chan *NetConn), maxIdleConn/2),
		maxOpenConn: maxOpenConn,
	}

	go pool.handleRequest()

	return pool
}

func (cp *ConnPoolImpl) handleRequest() {
	for {
		request, ok := <- cp.requestQueue
		if !ok {
			continue
		}

		conn, ok := <- cp.freeConns
		if !ok {
			break
		}

		select {
		case request <- conn:
		default:
			continue
		}
	}	
}

func (cp *ConnPoolImpl) GetConn(ctx context.Context) (*NetConn, error) {
	select{
	case <-ctx.Done():
		return nil, ERR_GET_CONN_TIMEOUT
	case conn, ok := <- cp.freeConns:
		if !ok {
			return nil, ERR_FREE_CONN_CHAN_CLOSED
		}
		return conn, nil
	default:
		return cp.requestNewConn(ctx)
	}
}

func (cp *ConnPoolImpl) requestNewConn(ctx context.Context) (*NetConn, error) {
	if cp.openConn < cp.maxOpenConn {
		return cp.createNewConn()
	}

	requester := make(chan *NetConn)
	defer close(requester)

	select{
	case <-ctx.Done():
		return nil, ERR_GET_CONN_TIMEOUT
	case cp.requestQueue <- requester:
		select {
		case conn, ok := <-requester:
			if !ok {
				return nil, ERR_FREE_CONN_CHAN_CLOSED
			}
			return conn, nil
		case <-ctx.Done():
			return nil, ERR_GET_CONN_TIMEOUT
		}
	default:
		return nil, CONN_POOL_QUEUE_FULL
	}
}

func (cp *ConnPoolImpl) createNewConn() (*NetConn, error) {
	cp.lock.Lock()
	cp.openConn += 1
	cp.lock.Unlock()

	return createNewConn()
}

func (cp *ConnPoolImpl) PutConn(ctx context.Context, conn *NetConn){
	select{
	case cp.freeConns <- conn:
		return
	default:
		cp.lock.Lock()
		cp.openConn -= 1
		cp.lock.Unlock()
		conn.Close()
		return
	}
}

var createNewConn = func() (*NetConn, error){
	return &NetConn{}, nil
}