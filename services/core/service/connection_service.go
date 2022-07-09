package service

type ConnectionService interface {
	Add(userID string, conn chan int)
	Ping(userID string)
}