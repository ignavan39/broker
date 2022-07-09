package connection

import "sync"

type ConnectionService struct {
	sustainableConnections map[string][]chan int  
	connectionLock sync.Mutex
}

func NewConnectionService() *ConnectionService {
	return &ConnectionService{
		sustainableConnections: make(map[string][]chan int),
	}
}

func (c *ConnectionService) Add(userID string, conn chan int) {
	c.connectionLock.Lock()
	defer c.connectionLock.Unlock()

	if _, ok := c.sustainableConnections[userID]; !ok {
		c.sustainableConnections[userID] = make([]chan int, 0)
	}

	c.sustainableConnections[userID] = append(c.sustainableConnections[userID], conn)
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