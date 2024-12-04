package db

import (
	users "db-worker/internal/models/request"
	"time"
)

type DBResult struct {
	ID             int64     `db:"id" json:"id"`
	OPERATION_ID   string    `db:"operation_id" json:"operation_id"`
	IN_PROGRESS    bool      `db:"in_progress" json:"in_progress"`
	DATA           []byte    `db:"data" json:"data"`
	OPERATION_TYPE string    `db:"type" json:"type"`
	CREATION_DATE  time.Time `db:"creation_date" json:"creation_date"`
	FINISH_DATE    time.Time `db:"finish_date" json:"finish_date"`
	VERSION        int64     `db:"version" json:"version"`
	USER_ID        int       `db:"user_id" json:"user_id"`
	FIRST_NAME     string    `db:"first_name" json:"first_name"`
	LAST_NAME      string    `db:"last_name" json:"last_name"`
	EMAIL          string    `db:"email" json:"email"`
}

type DBWorker interface {
	RegisterOperation(uniqID string, operation_type string, user_id int) error
	SetResult(uniqID string, data []byte) error
	GetResult(uniqID string) (dbResult DBResult, err error)
	GetAllOperations(limit int, operation_type string, user_id int) (dbOperations []DBResult, err error)
	GetOperation(uniqID string) (dbResult DBResult, err error)
	GetVersion(uniqID string) (version int64, err error)

	RegisterNewUser(email, password, firstName, lastName string) (int, error)
	CheckEmail(email string) (bool, error)
	CheckCorrectPassword(email, password string) (bool, int, error)
	GetUserByEmail(email string) (*users.DBUser, error)
	GetUserByID(id int) (*users.DBUser, error)
	GetAllUsers() ([]users.DBUser, error)
}
