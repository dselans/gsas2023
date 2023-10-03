package types

import (
	"encoding/json"
	"time"
)

const (
	SignupEventType EventType = iota
	LoginEventType
	LogoutEventType
	DeleteAccountEventType
)

type SignupEvent struct {
	Email         string    `json:"email"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	CountryCode   string    `json:"country_code"`
	SourceIP      string    `json:"source_ip"`
	Browser       string    `json:"browser"`
	OccurredAtUTC time.Time `json:"occurred_at_utc"`
}

type LoginEvent struct {
	Email         string    `json:"email"`
	Path          string    `json:"path"` // URL path that they used to login
	SourceIP      string    `json:"source_ip"`
	Browser       string    `json:"browser"`
	OccurredAtUTC time.Time `json:"occurred_at_utc"`
}

type EventType int

type Event struct {
	Type EventType
	Data json.RawMessage
}
