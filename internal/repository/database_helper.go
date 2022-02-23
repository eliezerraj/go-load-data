package repository

import (
	"context"
	"time"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/go-load-data/internal/core"

)

type DatabaseHelper interface {
	GetConnection(ctx context.Context) (*sql.DB, error)
	CloseConnection()
}

type DatabaseHelperImpl struct {
	client   	*sql.DB
	timeout		time.Duration
}

func NewDatabaseHelper(databaseRDS core.DatabaseRDS) (DatabaseHelper, error) {
	log.Printf("NewDatabaseHelper") 

	_ , cancel := context.WithTimeout(context.Background(), time.Duration(databaseRDS.Db_timeout)*time.Second)
	defer cancel()

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", databaseRDS.User, databaseRDS.Password, databaseRDS.Host, databaseRDS.Port, databaseRDS.DatabaseName) 
	
	fmt.Println("==========>", databaseRDS.Postgres_Driver, connStr)

	client, err := sql.Open(databaseRDS.Postgres_Driver, connStr)
	if err != nil {
		return DatabaseHelperImpl{}, err
	}
	err = client.Ping()
	if err != nil {
		return DatabaseHelperImpl{}, err
	}

	return DatabaseHelperImpl{
		client: client,
		timeout:  time.Duration(databaseRDS.Db_timeout) * time.Second,
	}, nil
}

func (d DatabaseHelperImpl) GetConnection(ctx context.Context) (*sql.DB, error) {
	//log.Printf("GetConnection") 
	return d.client, nil
}

func (d DatabaseHelperImpl) CloseConnection()  {
	log.Printf("CloseConnection !!!!") 
	defer d.client.Close()
}