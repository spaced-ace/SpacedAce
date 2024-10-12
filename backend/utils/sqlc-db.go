package utils

import (
	"fmt"
	"github.com/jackc/pgx/v5"
	"golang.org/x/net/context"
	"log"
	"os"
	"spaced-ace-backend/db"
	"sync"
	"time"
)

type SQLCQuerier struct {
	*db.Queries
	conn *pgx.Conn
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
		//host = "172.20.10.2"
		host = "192.168.1.72"
	}

	connectionString := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable", user, dbname, password, host, port)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, connectionString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	return &SQLCQuerier{
		Queries: db.New(conn),
		conn:    conn,
	}
}

func (q *SQLCQuerier) Close(ctx context.Context) error {
	mu.Lock()
	defer mu.Unlock()

	if q.conn != nil {
		err := q.conn.Close(ctx)
		if err != nil {
			return err
		}
		q.conn = nil
		instance = nil // Allow for re-initialization if needed
	}
	return nil
}

func (q *SQLCQuerier) IsConnected(ctx context.Context) bool {
	mu.Lock()
	defer mu.Unlock()

	if q.conn == nil {
		return false
	}

	return q.conn.Ping(ctx) == nil
}
