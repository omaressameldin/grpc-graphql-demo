package custom_models

import (
"time"
)

type Todo struct {
	ID     int
	Title   string
	Description   string
	IsDone   bool
	UserID int
	Reminder time.Time
}