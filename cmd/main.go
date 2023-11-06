package main

import (
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"io"
	"os"
	"webapp/config"
	"webapp/logger"
	"webapp/middleware"
	"webapp/models"
	"webapp/routes"
	"webapp/services"
)

func SetupGinRouter(services services.APIServices) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.NoRoute(routes.NoRouteHandler)
	r.Use(middleware.LogWebRequests())
	r.Any("/healthz", routes.HealthzGetReqHandler(services.Database))

	v1 := r.Group("/v1")
	{
		v1.Use(middleware.CheckDB(services), middleware.BasicAuth(services))
		v1.POST("/assignments", middleware.ValidateAssignmentsPayload(services), routes.AssignmentsPostHandler(services))
		v1.GET("/assignments/:id", routes.AssignmentGetByIDHandler(services))
		v1.GET("/assignments", routes.AssignmentGetHandler(services))
		v1.GET("/assignments/", routes.AssignmentGetHandler(services))
		v1.PUT("/assignments/:id", middleware.ValidateAssignmentsPayload(services), routes.AssignmentPutHandler(services))
		v1.DELETE("/assignments/:id", routes.AssignmentDeleteHandler(services))
		v1.PATCH("/assignments/:id", routes.AssignmentPatchHandler(services))
		v1.PATCH("/assignments/", routes.AssignmentPatchHandler(services))
	}

	return r
}

func init() {
	err := godotenv.Load()
	logger.InitLogger()
	logger.Info("Log location", zap.Any("log_location", os.Getenv("LOG_LOC")))
	logger.Info("Loading env file...")
	if err != nil {
		logger.Error("unable to load env file", zap.Any("error", err))
	} else {
		logger.Info("Loaded env file successfully")
	}
}

func loadDefaultAccounts(defaultUsers config.DefaultUsers, s services.APIServices) {
	if err := s.Database.Ping(); err != nil {
		logger.Warn("DB is unavailable, cannot load default users")
		return
	}

	file, err := os.Open(defaultUsers.Path)
	if err != nil {
		logger.Error("Error opening file", zap.Any("error", err))
		return
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Error("Error closing file", zap.Any("error", err))
		}
	}(file)

	csvReader := csv.NewReader(file)

	_, err = csvReader.Read()
	if err != nil {
		logger.Error("Error reading header", zap.Any("error", err))
		return
	}

	for {
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			logger.Error("Error reading record", zap.Any("error", err))
			return
		}

		passwordHash, hashError := s.AccountsService.HashPassword(record[3])
		if hashError != nil {
			logger.Error("Hashing password failed, %s\n", zap.Any("error", err))
		} else {
			account, accountFindError := s.AccountsService.GetAccountByEmail(record[2])

			if accountFindError != nil {
				logger.Warn("Account does not exists, %s\n", zap.Any("error", accountFindError))
				logger.Info("Creating new account")
				account = models.Account{
					FirstName: record[0],
					LastName:  record[1],
					Email:     record[2],
					Password:  passwordHash,
				}
				accountCreationError := s.AccountsService.AddAccount(account)
				if accountCreationError != nil {
					logger.Warn("Failed creating a default account, %s\n", zap.Any("error", accountCreationError))
					continue
				}
				logger.Info("Successfully created a default account", zap.Any("user", record[0]))
			} else {
				logger.Info("Account already exists", zap.Any("account", account))
			}
		}
	}
}

func main() {
	logger.Info("Configuring the application")
	configs := config.GetConfigs()

	logger.Info("Initializing application services")
	s := services.APIServices{}
	s.LoadServices(configs)

	logger.Info("Creating default users for the system")
	loadDefaultAccounts(configs.DefaultUsers, s)

	logger.Info("Setting up Gin logger")
	r := SetupGinRouter(s)

	logger.Info("Starting Gin webserver")
	err := r.Run(configs.ServerConfig.Host)
	if err != nil {
		logger.Fatal("Failed to start Gin server", zap.Error(err))
	}
}
