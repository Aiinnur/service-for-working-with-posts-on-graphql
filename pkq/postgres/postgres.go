package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"service-for-working-with-posts-on-graphql/internal/config"
	"service-for-working-with-posts-on-graphql/internal/repositories"
	"time"
)

func NewPostgresClient(cfg *config.Config) (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PostgreSQL config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	conn, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Проверка подключения
	if err = conn.Ping(ctx); err != nil {
		conn.Close()
		return nil, fmt.Errorf("could not acquire connection from PostgreSQL pool: %w", err)
	}

	_, err = conn.Exec(ctx, repositories.CreatePosts)
	if err != nil {
		return nil, fmt.Errorf("error creating posts table: %w", err)
	}

	_, err = conn.Exec(ctx, repositories.CreateComments)
	if err != nil {
		return nil, fmt.Errorf("error creating comments table: %w", err)
	}

	return conn, nil
}

func ClosePostgresClient(conn *pgxpool.Pool) {
	conn.Close()
}
