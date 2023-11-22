package respond

import "fmt"

// Response representing for http response data
type Response[T any] struct {
	status int
	data   T
}

// Error representing expected error from handler
type Error struct {
	Status  int
	Key     string
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("respond.Error{status:%d,key:%s,message:%s}", e.Status, e.Key, e.Message)
}

type Message struct {
	Name    string `json:"name"`
	Message string `json:"message,omitempty"`
}
