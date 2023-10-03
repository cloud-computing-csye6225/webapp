package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"io"
	"net/http"
	"webapp/config"
	"webapp/db"
)

func SetupGinRouter(db db.Database) *gin.Engine {
	r := gin.Default()

	r.Any("/healthz", func(context *gin.Context) {
		context.Header("Cache-Control", "no-store, no-cache, must-revalidate;")
		if context.Request.Method == http.MethodGet {
			querystring := context.Request.URL.RawQuery
			all, err := io.ReadAll(context.Request.Body)
			if err != nil {
				fmt.Printf("Error while reading the body %s\n", err)
			}
			if querystring != "" || len(all) > 0 {
				context.String(http.StatusBadRequest, "")
			} else {
				err := db.Ping()
				if err != nil {
					context.String(http.StatusServiceUnavailable, "")
					return
				}
				context.String(http.StatusOK, "")
			}
		} else {
			context.String(http.StatusMethodNotAllowed, "")

		}
	})

	r.NoRoute(func(context *gin.Context) {
		context.Header("Cache-Control", "no-store, no-cache, must-revalidate;")
		context.String(http.StatusNotFound, "")
	})
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
