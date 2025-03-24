package http

import "errors"

type Message struct {
	Data string `json:"message"`
}

func (m Message) Validate() error {
	if len(m.Data) == 0 {
		return errors.New("data not must be empty")
	}

	return nil
}