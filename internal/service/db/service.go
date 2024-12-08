package db

import (
	"db-worker/internal/config/db"
	model "db-worker/internal/models/db"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Worker struct {
	host            string
	port            string
	login           string
	password        string
	db_name         string
	table_name      string
	user_table_name string
	logger          *slog.Logger
}

func New(logger *slog.Logger) *Worker {
	return &Worker{
		host:            db.HOST,
		port:            db.PORT,
		login:           db.LOGIN,
		password:        db.PASSWORD,
		db_name:         db.DB_NAME,
		table_name:      db.RESULT_TABLE_NAME,
		user_table_name: db.USER_TABLE_NAME,
		logger:          logger,
	}
}

func (w *Worker) connectToDB() (*sqlx.DB, error) {
	connectionData := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s host=%s port=%s", w.login, w.db_name, w.password, w.host, w.port)
	return sqlx.Connect("postgres", connectionData)
}

func (w *Worker) RegisterOperation(uniqID string, operation_type string, user_id int) error {

	uniqID = strings.TrimSpace(uniqID)
	if len(uniqID) == 0 || len(uniqID) > 35 {
		return fmt.Errorf("uniqID is empty or too big")
	}

	db, err := w.connectToDB()
	if err != nil {
		w.logger.Error("Connection to DataBase", "error", err)
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO "+w.table_name+" (operation_id, in_progress, type, creation_date, finish_date, version, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7)", uniqID, true, operation_type, time.Now().Add(3*time.Hour), time.Now().Add(3*time.Hour), 0, user_id)
	if err != nil {
		w.logger.Error("Insert operation to DataBase", "error", err)
	}
	return err
}

func (w *Worker) SetResult(uniqID string, data []byte) error {

	uniqID = strings.TrimSpace(uniqID)
	if len(uniqID) == 0 || len(uniqID) > 35 {
		return fmt.Errorf("uniqID is empty or too big")
	}

	db, err := w.connectToDB()
	if err != nil {
		w.logger.Error("Connection to DataBase", "error", err)
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE "+w.table_name+" SET data = $1, in_progress = $2, finish_date = $3, version = version + 1 WHERE operation_id = $4", data, false, time.Now().Add(3*time.Hour), uniqID)
	if err != nil {
		w.logger.Error("Update operation to DataBase", "error", err)
	}
	return err
}

func (w *Worker) GetResult(uniqID string) (dbResult model.DBResult, err error) {

	uniqID = strings.TrimSpace(uniqID)
	if len(uniqID) == 0 || len(uniqID) > 35 {
		return model.DBResult{}, fmt.Errorf("uniqID is empty or too big")
	}

	db, err := w.connectToDB()
	if err != nil {
		w.logger.Error("Connection to DataBase", "error", err)
		return model.DBResult{}, err
	}
	defer db.Close()

	dbResults := make([]model.DBResult, 0, 2)

	err = db.Select(&dbResults, "SELECT * FROM "+w.table_name+" WHERE operation_id = $1", uniqID)
	if err != nil {
		w.logger.Error("Get operation from DataBase", "error", err)
		return model.DBResult{}, err
	}

	if len(dbResults) == 0 {
		return model.DBResult{}, fmt.Errorf("no results")
	}

	return dbResults[0], nil
}

func (w *Worker) GetAllOperations(limit int, operation_type string, user_id int) (dbResult []model.DBResult, err error) {
	w.logger.Info("GetAllOperations", "limit", limit, "operation_type", operation_type, "user_id", user_id)
	dbResult = make([]model.DBResult, 0, limit)
	db, err := w.connectToDB()
	if err != nil {
		w.logger.Error("Connection to DataBase", "error", err)
		return dbResult, err
	}
	defer db.Close()

	request := `SELECT o.id, o.operation_id, o.in_progress, o.type, o.creation_date, o.finish_date, o.version, o.user_id, o.data, u.first_name, u.last_name, u.email, u.user_status
	FROM result o
	JOIN users u ON o.user_id = u.id`

	if operation_type != "" {
		request += " WHERE o.type = '" + strings.ToLower(strings.TrimSpace(operation_type)) + "'"
	}

	if operation_type != "" && user_id >= 0 {
		request += " AND u.id = " + strconv.Itoa(user_id)
	} else if user_id >= 0 {
		request += " WHERE u.id = " + strconv.Itoa(user_id)
	}

	request += " ORDER BY o.id DESC"

	if limit > 0 {
		request += " LIMIT " + strconv.Itoa(limit)
	}

	err = db.Select(&dbResult, request)
	if err != nil {
		w.logger.Error("Get operation from DataBase", "error", err)
	}
	return dbResult, err
}

func (w *Worker) GetOperation(uniqID string) (dbResult model.DBResult, err error) {
	uniqID = strings.ToLower(strings.TrimSpace(uniqID))
	if len(uniqID) == 0 || len(uniqID) > 35 {
		return model.DBResult{}, fmt.Errorf("uniqID is empty or too big")
	}

	db, err := w.connectToDB()
	if err != nil {
		w.logger.Error("Connection to DataBase", "error", err)
		return model.DBResult{}, err
	}
	defer db.Close()
	q := `SELECT o.id, o.operation_id, o.in_progress, o.type, o.creation_date, o.finish_date, o.version, o.user_id, o.data, u.first_name, u.last_name, u.email 
	FROM result o
	JOIN users u ON o.user_id = u.id WHERE o.operation_id = $1 LIMIT 1`
	err = db.Get(&dbResult, q, uniqID)
	if err != nil {
		w.logger.Error("Get operation from DataBase", "error", err)
	}
	return dbResult, err
}

func (w *Worker) GetVersion(uniqID string) (version int64, err error) {

	type versionLabel struct {
		Version int64 `db:"version"`
	}

	uniqID = strings.TrimSpace(uniqID)
	if len(uniqID) == 0 || len(uniqID) > 35 {
		return 0, fmt.Errorf("uniqID is empty or too big")
	}

	db, err := w.connectToDB()
	if err != nil {
		w.logger.Error("Connection to DataBase", "error", err)
		return 0, err
	}

	var dbLabel = versionLabel{}

	defer db.Close()
	err = db.Get(&dbLabel, "SELECT version FROM "+w.table_name+" WHERE operation_id = $1 LIMIT 1", uniqID)
	return dbLabel.Version, err
}
