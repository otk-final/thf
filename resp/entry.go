package resp

import (
	"log"
)

var entryOK = "200"
var entryERR = "500"
var entryMessage = "success"

func InitRespEntryCode(ok, err string) {
	if ok == "" || err == "" {
		log.Fatalf("thf response entry code is illegal")
	}
	entryOK = ok
	entryERR = err
}
func InitRespEntryMessage(msg string) {
	entryMessage = msg
}

type Entry[T any] struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error"`
	Data    T      `json:"data"`
}

type Any struct {
	Entry[any]
}

func NewAny(data any) *Entry[any] {
	return NewEntry[any](data)
}

func NewEntry[T any](data T) *Entry[T] {
	return &Entry[T]{
		Data:    data,
		Code:    entryOK,
		Message: entryMessage,
	}
}

func NewError[T any](err error) *Entry[T] {
	rs := &Entry[T]{}
	rs.Code = entryERR
	if err != nil {
		rs.Error = err.Error()
	}
	return rs
}

func NewFail[T any](code string, err error) *Entry[T] {
	rs := &Entry[T]{}
	rs.Code = code
	if err != nil {
		rs.Error = err.Error()
	}
	return rs
}
