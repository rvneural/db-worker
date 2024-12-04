package db

import "os"

var (
	HOST     = os.Getenv("DB_HOST")
	PORT     = os.Getenv("DB_PORT")
	LOGIN    = os.Getenv("DB_LOGIN")
	PASSWORD = os.Getenv("DB_PASSWORD")

	DB_NAME = os.Getenv("DB_NAME")

	RESULT_TABLE_NAME = os.Getenv("RESULT_TABLE_NAME")
	USER_TABLE_NAME   = os.Getenv("USER_TABLE_NAME")
)
