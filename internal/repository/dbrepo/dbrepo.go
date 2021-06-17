package dbrepo

import (
	"database/sql"

	"github.com/andkolbe/go-websockets/internal/config"
	"github.com/andkolbe/go-websockets/internal/repository"
)

// only works with SQL
type mySQLDBRepo struct {
	App *config.AppConfig
	DB *sql.DB // database connection pool. *sql.DB is for mySQL
}

// same as above but for testing
type testDBRepo struct {
	App *config.AppConfig
	DB *sql.DB 
}

// lets us pass our connection pool and app config and return a repository
// because we return a pointer to mySQLDPRepo, it will connect to mySQL
// we can set up other repos for other dbs in the future
func NewMySQLRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &mySQLDBRepo {
		App: a,
		DB: conn,
	}
}

func NewTestingRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo {
		App: a,
	}
}