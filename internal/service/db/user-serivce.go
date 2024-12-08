package db

import (
	users "db-worker/internal/models/request"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func (w *Worker) RegisterNewUser(email, password, firstName, lastName string) (int, error) {
	db, err := w.connectToDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	query := `INSERT INTO users (email, password, first_name, last_name) VALUES ($1, $2, $3, $4) RETURNING id`
	password, err = w.hashPassword(password)
	if err != nil {
		return 0, err
	}
	var id int
	err = db.QueryRow(query, email, password, firstName, lastName).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (w *Worker) CheckEmail(email string) (bool, error) {
	db, err := w.connectToDB()
	if err != nil {
		return false, err
	}
	defer db.Close()

	query := `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err = db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (w *Worker) CheckCorrectPassword(email, password string) (bool, int, error) {
	db, err := w.connectToDB()
	if err != nil {
		return false, -1, err
	}
	defer db.Close()
	log.Println("Comparing for email:", email, "and password:", password)
	query := `SELECT password FROM users WHERE email = $1 LIMIT 1`
	var hashPass string
	var hashes []string
	err = db.Select(&hashes, query, email)
	if err != nil {
		log.Println("Selecting password error:", err)
		return false, -1, err
	}
	if len(hashes) == 0 {
		return false, -1, fmt.Errorf("Users not found")
	}
	hashPass = hashes[0]

	var correct = w.comparePassword(hashPass, password)
	var id = -1
	if correct {
		var ids []int
		query = `SELECT id FROM users WHERE email = $1 LIMIT 1`
		err = db.Select(&ids, query, email)
		if err != nil {
			return false, -1, err
		}
		if len(ids) == 0 {
			return false, -1, fmt.Errorf("Users not found")
		}
		id = ids[0]
	}

	return correct, id, nil
}

func (w *Worker) UpdatePassword(email, password string) error {
	db, err := w.connectToDB()
	if err != nil {
		return err
	}
	defer db.Close()

	query := `UPDATE users SET password = $1 WHERE email = $2`

	_, err = db.Exec(query, password, email)
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) GetUserByEmail(email string) (*users.DBUser, error) {
	db, err := w.connectToDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT id, email, first_name, last_name, user_status FROM users WHERE email = $1 LIMIT 1`

	var user users.DBUser
	err = db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.UserStatus)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (w *Worker) GetUserByID(id int) (*users.DBUser, error) {
	db, err := w.connectToDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT id, email, first_name, last_name, user_status FROM users WHERE id = $1 LIMIT 1`

	var user users.DBUser
	err = db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.UserStatus)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (w *Worker) hashPassword(password string) (string, error) {
	bytehash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytehash), err
}

func (w *Worker) comparePassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (w *Worker) GetAllUsers() ([]users.DBUser, error) {
	db, err := w.connectToDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT id, email, first_name, last_name FROM users`

	var users []users.DBUser
	err = db.Select(&users, query)
	if err != nil {
		return nil, err
	}

	return users, nil
}
