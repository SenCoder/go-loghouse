package response

import "net/http"

const (
	ErrOK         = iota + 200
	ErrBadRequest = http.StatusBadRequest
	ErrUnknown    = 499
)

type Response struct {
	Code uint32      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
