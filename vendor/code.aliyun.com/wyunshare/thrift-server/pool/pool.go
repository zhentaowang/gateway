package pool

import (
    "time"
    "sync"
    "container/list"
    "errors"
)

var (
    ErrPoolClosed    = errors.New("pool has been closed")
    ErrPoolExhausted = errors.New("pool max open client limit reached")
)

type Pool struct {
    // Dial is an application supplied function for creating new connections.
    Dial func() (interface{}, error)

    // Close is an application supplied functoin for closeing connections.
    Close func(c interface{}) error

    // TestOnBorrow is an optional application supplied function for checking
    // the health of an idle connection before the connection is used again by
    // the application. Argument t is the time that the connection was returned
    // to the pool. If the function returns an error, then the connection is
    // closed.
    TestOnBorrow func(c interface{}, t time.Time) error

    // Maximum number of idle connections in the pool.
    MaxIdle int

    // Maximum number of connections allocated by the pool at a given time.
    // When zero, there is no limit on the number of connections in the pool.
    MaxActive int

    // Close connections after remaining idle for this duration. If the value
    // is zero, then idle connections are not closed. Applications should set
    // the timeout to a value less than the server's timeout.
    IdleTimeout time.Duration

    // mu protects fields defined below.
    mu     sync.Mutex
    closed bool
    active int

    // Stack of idleConn with most recently used at the front.
    idle list.List
}

type idleConn struct {
    c interface{}
    t time.Time
}

// New creates a new pool. This function is deprecated. Applications should
// initialize the Pool fields directly as shown in example.
func New(dialFn func() (interface{}, error), closeFn func(c interface{}) error, maxIdle int) *Pool {
    return &Pool{Dial: dialFn, Close: closeFn, MaxIdle: maxIdle}
}

// Get gets a connection. The application must close the returned connection.
// This method always returns a valid connection so that applications can defer
// error handling to the first use of the connection.
func (p *Pool) Get() (interface{}, error) {
    p.mu.Lock()
    // if closed
    if p.closed {
        p.mu.Unlock()
        return nil, ErrPoolClosed
    }
    // Prune stale connections.
    if timeout := p.IdleTimeout; timeout > 0 {
        for i, n := 0, p.idle.Len(); i < n; i++ {
            e := p.idle.Back()
            if e == nil {
                break
            }
            ic := e.Value.(idleConn)
            if ic.t.Add(timeout).After(time.Now()) {
                break
            }
            p.idle.Remove(e)
            p.active -= 1
            p.mu.Unlock()
            // ic.c.Close()
            p.Close(ic.c)
            p.mu.Lock()
        }
    }
    // Get idle connection.
    if p.active >= p.MaxActive {
        p.mu.Unlock()
        return nil, ErrPoolExhausted
    }
    // No idle connection, create new.
    dial := p.Dial
    p.active += 1
    p.mu.Unlock()
    c, err := dial()
    if err != nil {
        p.mu.Lock()
        p.active -= 1
        p.mu.Unlock()
        c = nil
    }
    return c, err
}

// Put adds conn back to the pool, use forceClose to close the connection forcely
func (p *Pool) Put(c interface{}, forceClose bool) error {
    if !forceClose {
        p.mu.Lock()
        if !p.closed {
            p.idle.PushFront(idleConn{t: time.Now(), c: c})
            if p.idle.Len() > p.MaxIdle {
                // remove exceed conn
                c = p.idle.Remove(p.idle.Back()).(idleConn).c
            } else {
                c = nil
            }
        }
        p.mu.Unlock()
    }
    // close exceed conn
    if c != nil {
        p.mu.Lock()
        p.active -= 1
        p.mu.Unlock()
        return p.Close(c)
    }
    return nil
}

// ActiveCount returns the number of active connections in the pool.
func (p *Pool) ActiveCount() int {
    p.mu.Lock()
    active := p.active
    p.mu.Unlock()
    return active
}

// Relaase releases the resources used by the pool.
func (p *Pool) Release() error {
    p.mu.Lock()
    idle := p.idle
    p.idle.Init()
    p.closed = true
    p.active -= idle.Len()
    p.mu.Unlock()
    for e := idle.Front(); e != nil; e = e.Next() {
        p.Close(e.Value.(idleConn).c)
    }
    return nil
}
