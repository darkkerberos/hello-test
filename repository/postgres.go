package repository

import (
	"io"
)

type DBConfiguration struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBOptions  string
}

type PostgresReaderWriter interface {
	io.Closer
}
