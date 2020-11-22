package repository

import (
	"database/sql"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Repository interface {
	Store(id int64, username string, discriminator string) error
	Fetch(id int64) (username string, discriminator string, err error)
}

type repository struct {
	database *sql.DB
}

func (r *repository) Store(id int64, username string, discriminator string) error {
	statement, _ := r.database.Prepare("INSERT INTO users (id, username, discriminator) VALUES (?, ?, ?) ON CONFLICT (id) DO UPDATE SET username = ?, discriminator = ? WHERE id = ?")
	_, err := statement.Exec(id, username, discriminator, username, discriminator, id)
	return err
}

func (r *repository) Fetch(id int64) (username string, discriminator string, err error) {
	row := r.database.QueryRow("SELECT id, username, discriminator FROM users WHERE id = ?", id)
	err = row.Scan(&id, &username, &discriminator)
	return
}

func NewRepository() (*repository, error) {
	database, err := sql.Open("sqlite3", "./discord_users.db")
	if err != nil {
		return nil, err
	}
	_, err = database.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT, discriminator TEXT)")
	return &repository{database}, err
}

func ConvertStringToInt64(s string) (int64, error) {
	num, err := strconv.ParseUint(s, 10, 64)
	return int64(num), err
}
