package main

type Sender interface {
	Sender() string
	Receiver() string
	Send(msgSubject, msgText string) error
}
