package entity

import (
	"time"
)

type Task struct {
	Id           int64
	Name         string
	Text         string
	Time_add     time.Time
	Is_done      bool
	Time_done    *time.Time
	Is_important int
}

func NewTask(name, text string, isImp int) *Task {

	return &Task{
		Name:         name,
		Text:         text,
		Time_add:     time.Now(),
		Is_done:      false,
		Is_important: isImp,
	}
}

func (l *Task) DoTask() {
	l.Is_done = true
	tn := time.Now()
	l.Time_done = &tn
}
