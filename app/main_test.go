package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

func setupTestPool() *pgxpool.Pool {
	err := godotenv.Load("../.env")
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
	_, _ = pool.Exec(context.Background(), "TRUNCATE history RESTART IDENTITY CASCADE;")

	return pool
}

func TestCheckout(t *testing.T) {
	t.Run("Checkout", func(t *testing.T) {
		pool := setupTestPool()
		result, _ := Checkout(pool)
		require.NotNil(t, result)
		pool.Close()
	})
}

func TestGetHistory(t *testing.T) {
	t.Run("Checkout", func(t *testing.T) {
		pool := setupTestPool()
		_, _ = Checkout(pool)
		_, _ = Checkout(pool)

		result, _ := GetHistory(pool)
		require.Len(t, result, 2)
		require.Equal(t, int64(1), result[0].Id)
		pool.Close()
	})
}
