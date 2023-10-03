package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"webapp/config"
	"webapp/db"
	"webapp/routes"
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

func main() {
	configs := config.GetConfigs()

	d := &db.PostgresDB{}
	err := d.InitDatabase(configs.DBConfig)
	if err != nil {
		fmt.Printf("failed to initialize database: %s", err)
	}

	r := SetupGinRouter(d)

	err = r.Run(configs.ServerConfig.Host)
	if err != nil {
		panic("failed to start Gin server: " + err.Error())
	}
}
