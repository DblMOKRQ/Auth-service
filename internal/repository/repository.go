package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/DblMOKRQ/auth-service/internal/entity"
	"github.com/DblMOKRQ/auth-service/internal/token"

	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Register(user *entity.User) (int64, error) {
	psw, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {

		return 0, fmt.Errorf("failed to generate password hash: %v", err)
	}
	res, err := r.db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, psw)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.Code, sqlite3.ErrConstraint) {
				return 0, entity.ErrUserExists
			}
		}

		return 0, fmt.Errorf("repo.reg.failed to register user: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return id, nil
}

func (r *Repository) Login(user *entity.User, t *token.JWTMaker) (string, error) {
	res := r.db.QueryRow("SELECT username, password FROM users WHERE username = ?", user.Username)
	var username, psw string

	err := res.Scan(&username, &psw)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", entity.ErrUserNotFound
		}
		return "", fmt.Errorf("repo.login.failed to scan user: %w", err)
	}

	if username == "" {
		return "", entity.ErrUserNotFound
	}
	if psw == "" {
		return "", entity.ErrEmptyPassword
	}

	if err := bcrypt.CompareHashAndPassword([]byte(psw), []byte(user.Password)); err != nil {
		return "", entity.ErrInvalidPassword
	}

	tok, err := t.Create(user.Username)
	if err != nil {
		return "", fmt.Errorf("repo.login.failed to create token: %w", err)
	}

	return tok, nil
}

func (r *Repository) ValideToken(t *token.JWTMaker, token string) error {
	_, err := t.Validate(token)
	if err != nil {
		return fmt.Errorf("repo.validate.failed to validate token: %w", err)
	}
	return nil
}
