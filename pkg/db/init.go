package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

func InitConn() (*Queries, *pgxpool.Pool, error) {
	host := os.Getenv("DATABASE_HOST")
	if os.Getenv("TUNNEL_IP") != "" {
		host = os.Getenv("TUNNEL_IP")
	}
	username := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")
	connectstring := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable application_name=kidl-api", host, username, password, dbname)

	config, err := pgxpool.ParseConfig(connectstring)
	if err != nil {
		log.Panic("Error setting up pgxpool config", err)
	}

	config.MaxConns = 10
	config.MaxConnIdleTime = time.Second
	config.MinConns = 5
	config.HealthCheckPeriod = time.Second

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, nil, err
	}

	db := New(pool)

	return db, pool, nil
}
