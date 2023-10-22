package respond

import (
	"fmt"
	"net/http"
)

// Error representing expected error from handler
type Error struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
}

func (e Error) Error() string {
	return fmt.Sprintf("Error{code:%d,desc:%s}", e.Code, e.Desc)
}

var (
	ErrInternalServer = Error{Code: http.StatusInternalServerError, Desc: "Internal Server Error"}
)
