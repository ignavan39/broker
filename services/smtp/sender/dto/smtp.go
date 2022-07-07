package dto

type SendMailPayload struct {
	Subject string `json:"subject"`
	Message string `json:"message"`
	Recipient string `json:"recipient"`
	Sender string `json:"sender"`
}