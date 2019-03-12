package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	v1 "github.com/omaressameldin/grpc-graphql-demo/todos-server/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func BuildTaskResponse(t Task) (*v1.ToDo, error) {
	td := &v1.ToDo{
		Id: t.Key,
		Title: &v1.ToDo_TitleValue{
			TitleValue: t.Title,
		},
		Description: &v1.ToDo_DescriptionValue{
			DescriptionValue: t.Description,
		},
		IsDone: &v1.ToDo_IsDoneValue{
			IsDoneValue: t.IsDone,
		},
	}
	reminderValue, err := ptypes.TimestampProto(t.Reminder)
	if err != nil {
		return nil, status.Error(codes.Unknown, "reminder field has invalid format-> "+err.Error())
	}
	td.Reminder = &v1.ToDo_ReminderValue{
		ReminderValue: reminderValue,
	}
	return td, nil
}
