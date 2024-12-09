package utils

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/net/context"
	"log"
	"os"
	"spaced-ace-backend/db"
	"sync"
)

type SQLCQuerier struct {
	*db.Queries
	connPool *pgxpool.Pool
}

var (
	instance *SQLCQuerier
	once     sync.Once
	mu       sync.Mutex
)

func GetQuerier() *SQLCQuerier {
	once.Do(func() {
		instance = newSQLCQuerier()
	})
	return instance
}

func newSQLCQuerier() *SQLCQuerier {
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "postgres"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "test"
	}

	password := os.Getenv("DB_PASS")
	if password == "" {
		password = "test"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	connectionString := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable", user, dbname, password, host, port)

	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v", err)
	}

	// Create a connection pool
	connPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	return &SQLCQuerier{
		Queries:  db.New(connPool),
		connPool: connPool,
	}
}

func (q *SQLCQuerier) Close() error {
	mu.Lock()
	defer mu.Unlock()

	if q.connPool != nil {
		q.connPool.Close() // Close the connection pool
		q.connPool = nil
		instance = nil // Allow for re-initialization if needed
	}
	return nil
}

func (q *SQLCQuerier) IsConnected(ctx context.Context) bool {
	mu.Lock()
	defer mu.Unlock()

	if q.connPool == nil {
		return false
	}

	// Try acquiring a connection and immediately releasing it
	conn, err := q.connPool.Acquire(ctx)
	if err != nil {
		return false
	}
	defer conn.Release()

	return true
}
