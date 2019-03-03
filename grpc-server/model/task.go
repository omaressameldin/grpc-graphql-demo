package model

import (
	"encoding/json"
	"fmt"
	"time"
)

type Task struct {
	Key         int64     `json:"key"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Reminder    time.Time `json:"reminder"`
	IsDone      bool      `json:"isDone"`
}

func TaskToBuffer(t Task) ([]byte, error) {
	buf, err := json.Marshal(t)
	if err != nil {
		fmt.Println("marshal err:", err)
		return nil, err
	}
	return buf, nil
}

func BufferToTask(buf []byte) (Task, error) {
	var t Task
	var err error
	err = json.Unmarshal(buf, &t)
	if err != nil {
		fmt.Println("unmarshal err", err)
	}
	return t, err
}
