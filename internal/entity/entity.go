package entity

import (
	"time"
)

type Task struct {
	Name     string
	Text     string
	Time_add  time.Time
	Is_done     bool
	Time_done *time.Time
}

func NewTask(name, text string) *Task {
	return &Task{
		Name:    name,
		Text:    text,
		Time_add: time.Now(),
		Is_done:    false,
	}
}	

func (l *Task) DoTask() {
	l.Is_done = true
	tn := time.Now()
	l.Time_done = &tn
}
