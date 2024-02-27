package entity

import "fmt"

type GoDaddyError struct {
	StatusCode int
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (e GoDaddyError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
