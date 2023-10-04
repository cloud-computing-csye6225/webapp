package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"io"
	"os"
	"webapp/config"
	"webapp/db"
	"webapp/models"
	"webapp/routes"
	"webapp/services"
)

func SetupGinRouter(db db.Database) *gin.Engine {
	r := gin.Default()
	r.NoRoute(routes.NoRouteHandler)
	r.Any("/healthz", routes.HealthzGetReqHandler(db))
	return r
}

func init() {
	fmt.Printf("Reading configs from env...\n")
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("unable to load env file %s\n", err)
	} else {
		fmt.Print("Loaded env file successfully\n")
	}
}

func loadDefaultAccounts(defaultUsers config.DefaultUsers, s services.APIServices) {
	file, err := os.Open(defaultUsers.Path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(file)

	csvReader := csv.NewReader(file)

	_, err = csvReader.Read()
	if err != nil {
		fmt.Println("Error reading header:", err)
		return
	}

	for {
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading record:", err)
			return
		}

		// Process the CSV record (record is a []string)
		// You can access individual fields by index, e.g., record[0], record[1], etc.
		fmt.Println(record)
		passwordHash, hashError := s.AccountsService.HashPassword(record[3])
		if hashError != nil {
			fmt.Printf("Hashing password failed,\t%s\n", hashError)
		} else {
			account := models.Account{
				FirstName: record[0],
				LastName:  record[1],
				Email:     record[2],
				Password:  passwordHash,
			}
			accountCreationError := s.AccountsService.AddAccount(account)
			if accountCreationError != nil {
				fmt.Printf("Failed creating a default account,\t%s\n", accountCreationError)
				continue
			}
			fmt.Printf("Successfully created a default account for %s", record[0])
		}
	}
}

func main() {
	configs := config.GetConfigs()
	s := services.APIServices{}
	s.LoadServices(configs)
	//Load default accounts
	loadDefaultAccounts(configs.DefaultUsers, s)

	r := SetupGinRouter(s.Database)

	err := r.Run(configs.ServerConfig.Host)
	if err != nil {
		panic("failed to start Gin server: " + err.Error())
	}
}
