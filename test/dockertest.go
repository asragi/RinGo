package test

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
)

type PurgePool func() error

func CreateTestDB(
	imageName string,
	dockerfilePath string,
) (*sqlx.DB, PurgePool, error) {
	handleError := func(text string, err error) (*sqlx.DB, PurgePool, error) {
		return nil, nil, fmt.Errorf(text, err)
	}
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		return handleError("Could not construct pool: %w", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		return handleError("Could not connect to Docker: %w", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.BuildAndRun(imageName, dockerfilePath, []string{})
	if err != nil {
		return handleError("Could not start resource: %w", err)
	}
	/*
		if err = resource.Expire(20); err != nil {
			return handleError("Could not set expiration time: %w", err)
		}
	*/

	var db *sqlx.DB
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err = pool.Retry(
		func() error {
			db, err = sqlx.Open(
				"mysql",
				fmt.Sprintf("root:ringo@(localhost:%s)/mysql?parseTime=true", resource.GetPort("3306/tcp")),
			)
			if err != nil {
				return err
			}
			return db.Ping()
		},
	); err != nil {
		return handleError("Could not connect to database: %w", err)
	}

	stopPool := func() error {
		return pool.Purge(resource)
	}
	return db, stopPool, nil
}
