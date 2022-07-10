package connection

import (
	"sync"
)

type ConnectionService struct {
	sustainableConnections map[string]map[string]chan int
	connectionLock         sync.Mutex
}

func NewConnectionService() *ConnectionService {
	return &ConnectionService{
		sustainableConnections: make(map[string]map[string]chan int),
	}
}

func (c *ConnectionService) Add(userID string, key string, conn chan int) {
	c.connectionLock.Lock()
	defer c.connectionLock.Unlock()

	sc, ok := c.sustainableConnections[userID]
	if !ok {
		sc = make(map[string]chan int, 0)
		sc[key] = conn
		c.sustainableConnections[userID] = sc
	} else {
		c.sustainableConnections[userID][key] = conn
	}
}

func (c *ConnectionService) Ping(userID string) {
	c.connectionLock.Lock()
	defer c.connectionLock.Unlock()

	conns, ok := c.sustainableConnections[userID]

	if ok {
		for _, conn := range conns {
			conn <- 1
		}
	}
}

func (c *ConnectionService) Remove(key string) {
	c.connectionLock.Lock()
	defer c.connectionLock.Unlock()

	for userID, conns := range c.sustainableConnections {
		v, ok := conns[key]
		if ok {
			close(v)
			delete(conns, key)
			if len(conns) == 0 {
				delete(c.sustainableConnections, userID)
			} else {
				c.sustainableConnections[userID] = conns
			}
		}
	}
}
