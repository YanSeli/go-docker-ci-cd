package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"time"
)

type HistoryRecord struct {
	Id        int64     `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Print("Error loading .env file")
	}

	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbName := os.Getenv("POSTGRES_DATABASE")

	pool, err := pgxpool.Connect(context.Background(), fmt.Sprintf("postgresql://%s:%s@%s/%s", dbUser, dbPass, dbHost, dbName))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	r := gin.Default()

	r.POST("/", func(c *gin.Context) {
		var timestamp time.Time
		timestamp, err = Checkout(pool)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		c.JSON(201, gin.H{
			"timestamp": timestamp,
		})
	})

	r.GET("/", func(c *gin.Context) {
		var history []HistoryRecord
		history, err = GetHistory(pool)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		c.JSON(200, history)
	})

	err = r.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func Checkout(pool *pgxpool.Pool) (time.Time, error) {
	var timestamp = time.Now()
	_, err := pool.Exec(context.Background(), "insert into history(timestamp) values($1)", timestamp)
	if err != nil {
		return timestamp, err
	}
	return timestamp, nil

}

func GetHistory(pool *pgxpool.Pool) ([]HistoryRecord, error) {
	var history []HistoryRecord
	rows, err := pool.Query(context.Background(), "select id, timestamp from history")
	if err != nil {
		panic(errors.Wrap(err, "query bookings"))
	}

	for rows.Next() {
		h := HistoryRecord{}
		err := rows.Scan(&h.Id, &h.Timestamp)
		if err != nil {
			fmt.Println(err)
			continue
		}
		history = append(history, h)
	}

	return history, nil

}
