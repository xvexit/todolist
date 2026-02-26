package controllers

import (
	"encoding/json"
	"time"
)

type TaskDto struct {
	Name string
	Text string
}

type ErrorDto struct {
	Message string
	Time    time.Time
}

func NewErrorDto(mess string) *ErrorDto {
	return &ErrorDto{
		Message: mess,
		Time:    time.Now(),
	}
}

func (e *ErrorDto) ToString() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(b)
}
