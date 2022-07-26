package main

import (
    "fmt"
    "sync"
    "time"
	"log"
	"os"
	"strconv"
	"context"
	"net/http"
	"bytes"
	"encoding/json"

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

func worker_http(id string, bunch int,host string ,client http.Client) {
    fmt.Printf("worker_http %v %v starting\n", id, bunch)
	host = host + "/balance/save"
	for x := 1; x < bunch; x++{
		balance := NewBalance(x)
		payload := new(bytes.Buffer)
		json.NewEncoder(payload).Encode(balance)
		req_post , err := http.NewRequest("POST", host, payload)
		if err != nil {
			log.Println("Error http.NewRequest : ", err)
			panic(err)
		}

		req_post.Header = http.Header{
			"Accept_Language": []string{"pt-BR"},
			"jwt": []string{"cookie"},
			"Content-Type": []string{"application/json"},
		}
		resp, err := client.Do(req_post)
		if err != nil {
			log.Println("Error doing POST : ", err)
			panic(err)
		}
		defer resp.Body.Close()
		time.Sleep(time.Millisecond * time.Duration(1000))
	}

    fmt.Printf("worker_http %v done\n", id)
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
	log.Printf("------------------------")
	Configuration()
	log.Printf("------------------------")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Starting load-data 1.0")

	start := time.Now()
	grp := []string{"group-a"}

	if (rds_config.Type == "rest"){
		log.Printf("------------------------")
		log.Printf("REST (%v) ", rds_config.Url)
		log.Printf("-----------------------")
		client := http.Client{}

		var wg sync.WaitGroup
		wg.Add(len(grp))
		for _, f := range grp{
			go func(f string, b int, h string ,c http.Client){
				worker_http(f, b, h, c)
				wg.Done()
			}(f, rds_config.Bunch, rds_config.Url, client)
		}
		wg.Wait()

	}else {
		log.Printf("------------------------")
		log.Printf("DATABASE")
		log.Printf("-----------------------")
		repository_rds, err := repository.NewDatabaseHelper(rds_config)
		if err != nil {
			log.Print("Erro na abertura do Database", err)
			panic(err)
		}
		repo := repository.NewBalanceRepositoryRDS(repository_rds)
	
		var wg sync.WaitGroup
		wg.Add(len(grp))
	
		for _, f := range grp{
			go func(f string, b int, r repository.BalanceRepository){
				worker(f, b, r)
				wg.Done()
			}(f, rds_config.Bunch, repo)
		}
		wg.Wait()
		repository_rds.CloseConnection()
	}

	fmt.Printf("Time lapse : %s \n", time.Since(start))
}