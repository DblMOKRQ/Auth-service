package storage

import (
	"database/sql"
	"fmt"

	"github.com/DblMOKRQ/auth-service/internal/storage/sqlite"
)

func NewStorage(dbType string, storagePath string) (*sql.DB, error) {
	switch dbType {
	case "sqlite":
		return sqlite.NewStorage()

	default:
		return nil, fmt.Errorf("unknown storage type: %s", dbType)
	}
}
