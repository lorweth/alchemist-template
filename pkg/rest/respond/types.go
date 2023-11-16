package respond

import "fmt"

// Response representing for http response data
type Response[T any] struct {
	status int
	data   T
}

// Error representing expected error from handler
type Error struct {
	StatusCode int
	Name       string
	Message    string
}

func (e Error) Error() string {
	return fmt.Sprintf("respond.Error{code:%d,name:%s,msg:%s}", e.StatusCode, e.Name, e.Message)
}

type Message struct {
	Name    string `json:"name"`
	Message string `json:"message,omitempty"`
}
