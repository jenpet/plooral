package testutil

import (
	"database/sql"
	"fmt"
	"github.com/ory/dockertest/v3"
	log "github.com/sirupsen/logrus"
)

const (
	postgresHost               = "localhost"
	postgresUser               = "postgres"
	postgresPassword           = "postgres"

	postgresDockerImage    = "postgres"
	postgresDockerImageTag = "13"
)

var pool = newPool()

func newPool() DockerPool {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	return DockerPool{pool}
}

type DockerPool struct {
	pool *dockertest.Pool
}

func (dp *DockerPool) run(repository, tag string, env []string) *dockertest.Resource {
	res, err := dp.pool.Run(repository, tag, env)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	res.Expire(300)
	return res
}

func (dp *DockerPool) retry(op func() error) {
	if err := dp.pool.Retry(op); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
}

func Purge(name string) {
	if err := pool.pool.RemoveContainerByName(name); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func NewPostgres() (string,string){
	postgres := pool.run(postgresDockerImage, postgresDockerImageTag, []string{"POSTGRES_PASSWORD="+postgresPassword})
	uri := postgresConnectionURI(postgres)
	pool.retry(func() error {return postgresRetry(uri) })
	return postgres.Container.Name, uri
}

func postgresRetry(uri string) error {
	db, err := sql.Open("postgres", uri)
	if err != nil {
		return err
	}
	return db.Ping()
}

func postgresConnectionURI(r *dockertest.Resource) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		postgresUser,
		postgresPassword,
		postgresHost,
		r.GetPort("5432/tcp"),
		postgresUser)
}


