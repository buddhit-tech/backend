package config

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

func Connect() *pgxpool.Pool {
	globalEnv := GetEnv()
	postgresHost := globalEnv.PostgresHost
	postgresUser := globalEnv.PostgresUser
	postgresPassword := globalEnv.PostgresPassword
	postgresDb := globalEnv.PostgresDB
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s", postgresUser, postgresPassword, postgresHost, postgresDb)
	if dsn == "" {
		log.Fatal("DATABASE_URL env variable is required")
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("❌ Unable to connect to database: %v", err)
	}

	// check connection
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("❌ Cannot ping database: %v", err)
	}

	fmt.Println("✅ Connected to database")
	return pool
}
