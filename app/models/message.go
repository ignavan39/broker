package models

import "time"

type Message struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
	IsRead    bool       `json:"isRead"`
	Text      string     `json:"text"`
	User      User       `json:"user"`
	Peer      Peer       `json:"peer"`
	Parent    *Message   `json:"parent,omitempty"`
}
