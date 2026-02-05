package model

type Operation int

const (
	UnknownOperation Operation = iota
	SendMessageOperation
	UpdateMessageOperation
)
