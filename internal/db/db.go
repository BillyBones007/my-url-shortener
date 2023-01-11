package db

import (
	"context"

	"github.com/BillyBones007/my-url-shortener/internal/db/models"
	"github.com/BillyBones007/my-url-shortener/internal/hasher"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Интерфейс для работы с некой базой данных
type DBase interface {
	PostgresClient
	InsertURL(m *models.MainModel, h hasher.URLHasher) error
	SelectLongURL(m *models.Model) (*models.Model, error)
	UUIDIsExist(uuid string) bool
	SelectAllForUUID(uuid string) ([]models.Model, error)
}

// Интерфейс для работы с базой данных Postgresql через драйвер jackc/pgx
type PostgresClient interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	Close(ctx context.Context) error
	Ping(ctx context.Context) error
}
