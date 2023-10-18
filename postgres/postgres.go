package postgres

import (
	"database/sql"
	"fmt"
	"sync"

	// postgres
	"github.com/lib/pq"
	_ "github.com/lib/pq"

	"git.bluebird.id/logistic/commons/logger"
	repo "github.com/darkkerberos/hello-test/repository"
	"go.elastic.co/apm/module/apmsql"
)

var mutex = &sync.RWMutex{}

// ReaderWriter ...
type ReaderWriter struct {
	db    *sql.DB
	mutex sync.RWMutex
}

const driverName = "postgres"

// NewPostgresReaderWriter ...
func NewPostgresReaderWriter(conf repo.DBConfiguration) (repo.PostgresReaderWriter, error) {
	apmsql.Register(driverName, &pq.Driver{})
	connURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s%s",
		conf.DBUser,
		conf.DBPassword,
		conf.DBHost,
		conf.DBPort,
		conf.DBName,
		conf.DBOptions)

	logger.Debug(fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s%s",
		conf.DBUser,
		"*************",
		conf.DBHost,
		conf.DBPort,
		conf.DBName,
		conf.DBOptions))

	db, err := apmsql.Open(driverName, connURL)
	if err != nil {
		return nil, err
	}

	return &ReaderWriter{
		db:    db,
		mutex: sync.RWMutex{},
	}, nil
}

// Close ...
func (p *ReaderWriter) Close() error {
	return p.db.Close()
}
