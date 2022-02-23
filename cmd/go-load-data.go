package main

import (
    "fmt"
    "sync"
    "time"
	"log"
	"os"
	"strconv"
	"context"

	"github.com/spf13/viper"

	"github.com/go-load-data/internal/core"
	"github.com/go-load-data/internal/repository"

)

var (
	rds_config 			core.DatabaseRDS
	_applicationFileName = "application"
	_extension           = "yaml"
	_resourcePath        = "../resources/"
)

func Configuration() {
	fmt.Println("Configuration " ,_resourcePath , _applicationFileName,_extension)
	
	viper.AddConfigPath(_resourcePath)
	viper.SetConfigName(_applicationFileName)
	viper.SetConfigType(_extension)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Not Found : ",err)
			os.Exit(1)
		} else {
			log.Println("Erro Fatal : ",err)
			os.Exit(1)
		}
	}
	
	err := viper.Unmarshal(&rds_config)
	if err != nil {
		log.Println("Erro Unmarshall : ",err)
		os.Exit(1)
	}
	log.Println(rds_config)
}

func worker(id string, bunch int, repo repository.BalanceRepository) {
    fmt.Printf("Worker %v %v starting\n", id, bunch)

	for x := 1; x < bunch; x++{
		balance := NewBalance(x)
		_ , err := repo.Save(context.Background(), balance)
		if err != nil {
			log.Print("Erro no Save", err)
			panic(err)
		}
	}

    fmt.Printf("Worker %v done\n", id)
}

func NewBalance(i int) core.Balance{
	acc := "acc-" + strconv.Itoa(i)
	description := "COOKIE-"+ strconv.Itoa(i) + " - OK"
	
	balance := core.Balance{
		Id:    int32(i),
		Account: acc,
		Amount: 1,
		DateBalance: time.Now(),
		Description: description,
	}
	return balance
}

func init() {
	fmt.Println("init")
	Configuration()
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Starting load-data 1.0")

	start := time.Now()
	foo := []string{"group-a","group-b","group-c","group-d","group-e","group-f","group-g","group-h","group-i","group-j"}

	repository_rds, err := repository.NewDatabaseHelper(rds_config)
	if err != nil {
		log.Print("Erro na abertura do Database", err)
		panic(err)
	}
	repo := repository.NewBalanceRepositoryRDS(repository_rds)

    var wg sync.WaitGroup
	wg.Add(len(foo))

	for _, f := range foo{
		go func(f string, b int, r repository.BalanceRepository){
			worker(f, b, r)
			wg.Done()
		}(f, rds_config.Bunch, repo)
	}
	wg.Wait()
	repository_rds.CloseConnection()
	fmt.Printf("Time lapse : %s \n", time.Since(start))
}