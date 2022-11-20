package models

type Action struct {
	Key     string `json:"key"`
	Value   []byte `json:"value"`
	Command string `json:"command"`
}

const (
	GET    = "GET"
	SET    = "SET"
	UPDATE = "UPDATE"
	DELETE = "DELETE"
	ERROR  = "ERROR"
)
