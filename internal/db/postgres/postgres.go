package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/BillyBones007/my-url-shortener/internal/db/models"
	"github.com/BillyBones007/my-url-shortener/internal/hasher"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Тип для работы с Postgresql
type ClientPostgres struct {
	Pool     *pgxpool.Pool
	ConfigCP *pgxpool.Config
}

// Конструктор нового соединения с базой Postgresql
func NewClientPostgres(dst string) *ClientPostgres {
	config, err := pgxpool.ParseConfig(dst)
	if err != nil {
		log.Fatal(err)
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}
	cp := ClientPostgres{Pool: pool, ConfigCP: config}
	return &cp
}

func (c *ClientPostgres) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	ct := pgconn.CommandTag{}
	return ct, nil
}

func (c *ClientPostgres) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return nil, nil
}

func (c *ClientPostgres) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return nil
}

func (c *ClientPostgres) Begin(ctx context.Context) (pgx.Tx, error) {
	return nil, nil
}

func (c *ClientPostgres) Close(ctx context.Context) error {
	return nil
}

func (c *ClientPostgres) Ping(ctx context.Context) error {
	if err := c.Pool.Ping(ctx); err != nil {
		return fmt.Errorf("ERROR: couldn't ping postgres database: %v", err)
	}
	return nil
}

func (c *ClientPostgres) InsertURL(m *models.MainModel, h hasher.URLHasher) error {
	return nil
}

func (c *ClientPostgres) SelectLongURL(m *models.Model) (*models.Model, error) {
	return nil, nil
}

func (c *ClientPostgres) UUIDIsExist(uuid string) bool {
	return false
}

func (c *ClientPostgres) SelectAllForUUID(uuid string) ([]models.Model, error) {
	return nil, nil
}
