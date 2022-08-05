package core

import (
	"time"

)

type Balance struct {
	Id					string		`json:"id"`
	BalanceId			string		`json:"balance_id"`
    Account 			string 		`json:"account"`
	Amount				int32 		`json:"amount"`
    DateBalance  		time.Time 	`json:"date_balance"`
	Description			string 		`json:"description"`
}

type DatabaseRDS struct {
    Host 				string `json:"host"`
    Port  				string `json:"port"`
	Schema				string `json:"schema"`
	DatabaseName		string `json:"databaseName"`
	User				string `json:"user"`
	Password			string `json:"password"`
	Db_timeout			int		`json:"db_timeout"`
	Postgres_Driver		string `json:"postgres_driver"`
	Bunch				int		`json:"bunch"`
	Url					string `json:"url"`
	Type				string `json:"type"`
}