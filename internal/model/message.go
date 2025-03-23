package model

type Message []byte

func (m Message) String() string {return string(m)}

type QueueID string