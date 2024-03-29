package connpool

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/leminhviett/TCP-server/config"
	"github.com/leminhviett/TCP-server/domain/customError"
)

type NetConn struct {
	net.Conn
}

type ConnPool interface {
	GetConn(ctx context.Context) (*NetConn, error)
	PutConn(ctx context.Context, conn *NetConn)
}

type ConnPoolImpl struct {
	freeConn       chan *NetConn
	requestQueue   chan (chan *NetConn)
	maxOpenConn    int32
	numberOpenConn int32
	lock           sync.Mutex
}

func NewConnPool(ctx context.Context, maxIdleConn, maxOpenConn int32) *ConnPoolImpl {
	pool := &ConnPoolImpl{
		freeConn:     make(chan *NetConn, maxIdleConn),
		requestQueue: make(chan (chan *NetConn), maxOpenConn/2),
		maxOpenConn:  maxOpenConn,
	}

	fmt.Printf("freeConn cap: %d, queue cap: %d, maxOpenConn: %d \n", cap(pool.freeConn), cap(pool.requestQueue), maxOpenConn)

	go pool.requestHandler(ctx)

	return pool
}

func (cp *ConnPoolImpl) requestHandler(ctx context.Context) {
	for {
		request, ok := <-cp.requestQueue
		if !ok {
			break //channel is closed
		}

		conn, ok := <-cp.freeConn
		if !ok {
			break //channel is closed
		}

		select {
		case request <- conn:
		default:
			cp.PutConn(ctx, conn) // handle edge case when requester close the channel
			continue
		}
	}
}

func (cp *ConnPoolImpl) GetConn(ctx context.Context) (*NetConn, error) {
	select {
	case <-ctx.Done():
		return nil, customError.ERR_GET_CONN_TIMEOUT
	case conn, ok := <-cp.freeConn:
		if !ok {
			return nil, customError.ERR_FREE_CONN_CHAN_CLOSED
		}
		return conn, nil
	default:
		return cp.requestNewConn(ctx)
	}
}

func (cp *ConnPoolImpl) requestNewConn(ctx context.Context) (*NetConn, error) {
	cp.lock.Lock()
	if cp.maxOpenConn > cp.numberOpenConn {
		conn, err := createNewConn()
		if err != nil {
			cp.lock.Unlock()
			return nil, err
		}
		cp.numberOpenConn += 1
		cp.lock.Unlock()
		return conn, nil
	}
	cp.lock.Unlock()

	requester := make(chan *NetConn)

	select {
	case <-ctx.Done():
		return nil, customError.ERR_GET_CONN_TIMEOUT
	case cp.requestQueue <- requester:
		select {
		case conn := <-requester:
			return conn, nil
		case <-ctx.Done():
			return nil, customError.ERR_GET_CONN_TIMEOUT
		}
	}
}

func (cp *ConnPoolImpl) PutConn(_ context.Context, conn *NetConn) {
	select {
	case cp.freeConn <- conn:
	default:
		cp.lock.Lock()
		cp.numberOpenConn -= 1
		cp.lock.Unlock()
		err := conn.Close()
		if err != nil {
			return
		}
	}
}

var createNewConn = func() (*NetConn, error) {
	conn, err := net.Dial(config.TCP_CONN_TYPE,
		fmt.Sprintf("%s:%s", config.TCP_SERVER_CONN_HOST, config.TCP_SERVER_CONN_PORT))

	return &NetConn{conn}, err
}
