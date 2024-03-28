package repository

import (
	"database/sql"
	"lab6_2/internal/model"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(connection string) (*Repository, error) {
	db, err := sql.Open("mysql", connection)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) Close() error {
	err := r.db.Close()
	return err
}

func (r *Repository) GetUsers() ([]*model.User, error) {
	q := "SELECT username, email, message FROM iu9budnikov"
	rows, err := r.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]*model.User, 0)
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.Name, &u.Email, &u.Message); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}

func (r *Repository) SaveLog(email, message string, status bool) error {
	q, err := r.db.Prepare(
		"INSERT INTO iu9budnikov_logs (email, with_error, message, time) VALUES (?, ?, ?, ?)",
	)
	if err != nil {
		return err
	}
	defer q.Close()
	_, err = q.Exec(email, status, message, time.Now().UTC())
	if err != nil {
		return err
	}

	qTime, err := r.db.Prepare(
		"INSERT INTO iu9budnikov (time) VALUES (?)",
	)
	if err != nil {
		return err
	}
	defer qTime.Close()
	_, err = qTime.Exec(time.Now().UTC())
	if err != nil {
		return err
	}

	return nil
}
