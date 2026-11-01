package main

import (
	"fmt"
	"log"
	"medicity/config"
	"medicity/database"
	"medicity/logger"

	//"medicity/pkg/utils"
	"medicity/internal/routes"
	//"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {

	confg := config.LoadConfig()
	if err := logger.InitLogger(); err != nil {
		log.Fatal(err)
	}

	defer logger.Sync()

	logger.Log.Info("Logger initialized")
	//connect DB
	database.ConnectDatabase(confg)

	// Run migrations
	if err := database.MigrateDatabase(); err != nil {
		panic(err)
	}

	r := gin.New()
	r.Use(ginzap.RecoveryWithZap(logger.Log, true))

	//load HTML templates
	r.LoadHTMLGlob("../templates/*")
	//load static files
	r.Static("/static", "../static")

	 store := cookie.NewStore([]byte("your-secret-key"))
    r.Use(sessions.Sessions("mysession", store))

	routes.SetupRoutes(r)

	logger.Log.Info("Server started", zap.String("port", confg.Port))
	//Run server
	fmt.Println("Server started on port", confg.Port)
	r.Run(":" + confg.Port)

}
