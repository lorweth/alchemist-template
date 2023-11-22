package httpio

import (
	"fmt"
)

// Error represents an application-specific error structure.
// It implements the error interface, allowing it to be used as an error type.
type Error struct {
	Status int
	Key    string
	Desc   string
}

func (e Error) Error() string {
	return fmt.Sprintf("respond.Error{status:%d,key:%s,desc:%s}", e.Status, e.Key, e.Desc)
}

// Message represents a generic message structure used for communication.
type Message struct {
	Key     string `json:"key"`
	Content string `json:"content,omitempty"`
}
