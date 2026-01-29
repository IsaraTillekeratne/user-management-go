package testutils

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresEnv struct {
	DB        *sql.DB
	Terminate func()
}

func SetUpPostgresEnv(ctx context.Context) *PostgresEnv {

	// temporary postgres container for tests
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "useradmin",
			"POSTGRES_PASSWORD": "userpassword",
			"POSTGRES_DB":       "userdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err)
	}

	host, err := postgres.Host(ctx)
	if err != nil {
		log.Fatal(err)
	}

	port, err := postgres.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=useradmin password=userpassword dbname=userdb sslmode=disable", host, port.Port())
	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return &PostgresEnv{
		DB: dbConn,
		Terminate: func() {
			if err := postgres.Terminate(ctx); err != nil {
				log.Printf("failed to terminate postgres container: %v", err)
			}
		},
	}
}
