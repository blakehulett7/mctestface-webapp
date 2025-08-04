package dbrepo

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "postgres"
	dbName   = "users_test"
	port     = "5435"
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5"
)

var resource *dockertest.Resource
var pool *dockertest.Pool
var testDB *sql.DB

func TestMain(m *testing.M) {
	log.Println("Sancta Maria...")
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Docker not running: %s\n", err)
	}

	pool = p
	options := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14.5",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	resource, err := pool.RunWithOptions(&options)
	if err != nil {
		pool.Purge(resource)
		log.Fatalf("Could not start resource: %s\n", err)
	}

	err = pool.Retry(func() error {
		var innerErr error
		testDB, innerErr = sql.Open("pgx", fmt.Sprintf(dsn, host, port, user, password, dbName))
		if innerErr != nil {
			log.Println("Error:", innerErr)
			return innerErr
		}

		return testDB.Ping()
	})
	if err != nil {
		pool.Purge(resource)
		log.Fatalf("Could not connect to db: %s\n", err)
	}

	err = CreateTables()
	if err != nil {
		log.Fatalf("Could not populate db: %s\n", err)
	}

	code := m.Run()

	err = pool.Purge(resource)
	if err != nil {
		log.Fatalf("could not purge resouces: %s\n", err)
	}
	os.Exit(code)
}

func CreateTables() error {
	tableSql, err := os.ReadFile("./testdata/users.sql")
	if err != nil {
		return err
	}

	_, err = testDB.Exec(string(tableSql))
	if err != nil {
		return err
	}

	return nil
}

func TestPingDb(t *testing.T) {
	err := testDB.Ping()
	if err != nil {
		t.Error("can't ping db...")
	}
}
