package repository

import (
	"context"
	"log"
	"time"

	"github.com/go-load-data/internal/core"

)

type BalanceRepository interface {
	Save(ctx context.Context, balance core.Balance) (core.Balance, error)
}

type BalanceRepositoryRDSImpl struct {
	DatabaseHelper DatabaseHelper
}

func NewBalanceRepositoryRDS(databaseHelper DatabaseHelper) BalanceRepository {
	return BalanceRepositoryRDSImpl{
		DatabaseHelper: databaseHelper,
	}
}

func (b BalanceRepositoryRDSImpl) Save(ctx context.Context, balance core.Balance) (core.Balance, error) {
	//log.Printf("Save") 

	client, _ := b.DatabaseHelper.GetConnection(ctx)

	stmt, err := client.Prepare(`INSERT INTO balance ( balance_id, 
														 account, 
														 amount, 
														 date_balance, 
														 Description) 
									VALUES( $1, $2, $3, $4, $5) `)
	if err != nil {
		log.Printf("Error Database - 1 : ", err) 
		return core.Balance{}, err
	}
	//log.Printf("Save 2 ", balance) 
	_, err = stmt.Exec(	balance.Id, 
						balance.Account,
						balance.Amount,
						time.Now(),
						balance.Description)
 	if err != nil {
		log.Printf("Error Database - 2 : ", err) 
		return core.Balance{}, err
	}
	return balance , nil
}