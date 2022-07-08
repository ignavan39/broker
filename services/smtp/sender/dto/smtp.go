package dto

type SendMailPayload struct {
	Subject     string             `json:"subject"`
	Message     string             `json:"message"`
	Recipient   string             `json:"recipient"`
	Attachments *map[string]string `json:"attachments,omitempty"`
}
