package main

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

/*
	My experiment implementing a connection pool
	@author: raman-rai
*/

type ConnectionPool struct {
	connections []*sql.DB
	channel     chan interface{}
	mu          sync.Mutex
}

func (pool *ConnectionPool) Get() *sql.DB {
	<-pool.channel

	pool.mu.Lock()
	defer pool.mu.Unlock()

	conn := pool.connections[0]
	pool.connections = pool.connections[1:]

	return conn
}

func (pool *ConnectionPool) Put(conn *sql.DB) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	pool.connections = append(pool.connections, conn)

	pool.channel <- nil
}

func NewPool(maxSize int) (*ConnectionPool, error) {
	pool := &ConnectionPool{
		connections: make([]*sql.DB, 0, maxSize),
		channel:     make(chan interface{}, maxSize),
	}

	for i := 0; i < maxSize; i++ {
		conn, err := NewDBConnection()
		if err != nil {
			fmt.Println("Error creating a new database connection:", err)
			pool.Close()
			return nil, err
		}
		pool.connections = append(pool.connections, conn)
		pool.channel <- nil
	}

	return pool, nil
}

func (pool *ConnectionPool) Close() {
	for _, conn := range pool.connections {
		conn.Close()
	}
}

func NewDBConnection() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		panic(err)
	}

	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("MYSQL_USER")
	cfg.Passwd = os.Getenv("MYSQL_PASSWORD")
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "testdb"

	return sql.Open("mysql", cfg.FormatDSN())
}

func QueryDB(conn *sql.DB, query string) (sql.Result, error) {
	return conn.Exec(query)
}

func ExecuteQueriesWithPool() {
	connectionPool, err := NewPool(10)
	if err != nil {
		fmt.Println("Error initializing connection pool:", err)
		return
	}
	defer connectionPool.Close()

	fmt.Println("Connections in connection pool:", len(connectionPool.connections))
	var wg sync.WaitGroup

	start := time.Now()
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			conn := connectionPool.Get()

			_, err = QueryDB(conn, "SELECT SLEEP(0.01)")
			if err != nil {
				fmt.Println("Error executing query:", err)
			}

			connectionPool.Put(conn)
		}()
	}

	wg.Wait()

	fmt.Println("Total time taken:", time.Since(start))
}

func main() {
	ExecuteQueriesWithPool()
}
