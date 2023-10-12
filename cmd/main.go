package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"io"
	"os"
	"webapp/config"
	"webapp/middleware"
	"webapp/models"
	"webapp/routes"
	"webapp/services"
)

func SetupGinRouter(services services.APIServices) *gin.Engine {
	r := gin.Default()
	r.NoRoute(routes.NoRouteHandler)
	r.Any("/healthz", routes.HealthzGetReqHandler(services.Database))

	v1 := r.Group("/v1")
	{
		v1.POST("/assignments", middleware.BasicAuth(services), middleware.ValidateAssignmentsPayload(services), routes.AssignmentsPostHandler(services))
		v1.GET("/assignments/:id", middleware.BasicAuth(services), routes.AssignmentGetByIDHandler(services))
		v1.GET("/assignments", middleware.BasicAuth(services), routes.AssignmentGetHandler(services))
		v1.GET("/assignments/", middleware.BasicAuth(services), routes.AssignmentGetHandler(services))
		v1.PUT("/assignments/:id", middleware.BasicAuth(services), middleware.ValidateAssignmentsPayload(services), routes.AssignmentPutHandler(services))
		v1.DELETE("/assignments/:id", middleware.BasicAuth(services), routes.AssignmentDeleteHandler(services))
		v1.PATCH("/assignments/:id", middleware.BasicAuth(services), routes.AssignmentPatchHandler(services))
		v1.PATCH("/assignments/", middleware.BasicAuth(services), routes.AssignmentPatchHandler(services))
	}

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

		passwordHash, hashError := s.AccountsService.HashPassword(record[3])
		if hashError != nil {
			fmt.Printf("Hashing password failed,\t%s\n", hashError)
		} else {
			account, accountFindError := s.AccountsService.GetAccountByEmail(record[2])

			if accountFindError != nil {
				fmt.Printf("Account does not exists,\t%s\n", accountFindError)
				fmt.Println("Creating new account")
				account = models.Account{
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
				fmt.Printf("Successfully created a default account for %s\n", record[0])
			} else {
				fmt.Printf("Account already exists, \t%s\n", account)
			}
		}
	}
}

func main() {
	configs := config.GetConfigs()
	s := services.APIServices{}
	s.LoadServices(configs)
	//Load default accounts
	loadDefaultAccounts(configs.DefaultUsers, s)

	r := SetupGinRouter(s)

	err := r.Run(configs.ServerConfig.Host)
	if err != nil {
		panic("failed to start Gin server: " + err.Error())
	}
}
