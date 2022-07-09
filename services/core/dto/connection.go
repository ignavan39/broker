package dto

type Meta struct {
	QueueName    string `json:"queueName"`
	ExchangeName string `json:"exchange"`
}

type CreateConnectionBase struct {
	Meta
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Vhost    string `json:"vhost"`
	Password string `json:"password"`
}
