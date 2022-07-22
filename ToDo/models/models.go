package models

import "time"

type Status string

const (
	StatusDraft  = "draft"
	StatusActive = "active"
)

type AddLogin struct {
	ID       string `db:"id" json:"ID"`
	Password string `db:"password" json:"Password"`
}

type AddTask struct {
	Task   string `db:"task" json:"Task"`
	Status Status `db:"status" json:"Status"`
}

type CreateUser struct {
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type DeleteTask struct {
	TaskID string `db:"task_id" json:"TaskID"`
}

type Session struct {
	SessionId string    `db:"session_id " json:"SessionId"`
	ExpiredAt time.Time `db:"expired_at " json:"ExpiredAt"`
	ID        string    `db:"id" json:"ID"`
}

type Task struct {
	Task string `db:"task" json:"Task"`
}

type UpdateTask struct {
	TaskID string `db:"task_id" json:"TaskID"`
	Task   string `db:"task" json:"Task"`
}

type User struct {
	TaskID string `db:"task_id" json:"TaskID"`
	Task   string `db:"task" json:"Task"`
	Status Status `db:"status" json:"Status"`
}

type UserId struct {
	ID string `json:"id" db:"ID"`
}

type PaginatedData struct {
	TotalCount int    `json:"-" db:"total_count"`
	TaskID     string `db:"task_id" json:"TaskID"`
	Task       string `db:"task" json:"Task"`
	Status     Status `db:"status" json:"Status"`
}

type PaginatedTask struct {
	Data       []PaginatedData `json:"data"`
	TotalCount int             `json:"totalCount"`
}
