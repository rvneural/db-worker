package db

import "time"

type DBResult struct {
	ID             int64     `db:"id"`
	OPERATION_ID   string    `db:"operation_id"`
	IN_PROGRESS    bool      `db:"in_progress"`
	DATA           []byte    `db:"data"`
	Error          string    `db:"error"`
	OPERATION_TYPE string    `db:"type"`
	CREATION_DATE  time.Time `db:"creation_date"`
	FINISH_DATE    time.Time `db:"finish_date"`
	VERSION        int64     `db:"version"`
}

type DBWorker interface {
	RegisterOperation(uniqID string, operation_type string) error
	SetResult(uniqID string, data []byte) error
	GetResult(uniqID string) (dbResult DBResult, err error)
	GetAllOperations(limit int, operation_type string) (dbOperations []DBResult, err error)
	GetOperation(uniqID string) (dbResult DBResult, err error)
	UpdateResult(uniqID string, data []byte) (err error)
	GetVersion(uniqID string) (version int64, err error)
}
