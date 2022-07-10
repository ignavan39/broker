package service

type ConnectionService interface {
	Add(userID string, key string, conn chan int)
	Ping(userID string)
	Remove(key string)
}
