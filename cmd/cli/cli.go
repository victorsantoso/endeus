package cli

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/victorsantoso/endeus/internal"

	userHandler "github.com/victorsantoso/endeus/users/http/handler"
	userRepository "github.com/victorsantoso/endeus/users/repository"
	userUsecase "github.com/victorsantoso/endeus/users/usecase"
	authMiddleware "github.com/victorsantoso/endeus/users/http/middleware"

	recipeHandler "github.com/victorsantoso/endeus/recipes/http/handler"
	recipeRepository "github.com/victorsantoso/endeus/recipes/repository"
	recipeUsecase "github.com/victorsantoso/endeus/recipes/usecase"

)

func Bootstrap(configPath string) error {
	defer func() {
		if err := recover(); err != nil {
			log.Warn("panic occured")
		}
	}()
	g := gin.New()
	g.Use(gin.Recovery())
	// configure application
	application := internal.ConfigureApplication()
	if application.Debug {
		internal.SetDebug(application.Debug)
	}
	// configure db config
	db := internal.ConfigureDatabase()
	// configure db connection
	dbConn := internal.NewPostgresConn(db)
	// define domain of applications
	// user domain
	userRepository := userRepository.NewUserRepository(dbConn)
	userUsecase := userUsecase.NewUserUsecase(userRepository)
	userHandler.NewUserHandler(g, userUsecase)
	// set authentication middleware
	authMiddleware := authMiddleware.AuthMiddleware(userRepository)
	// recipe domain
	recipeRepository := recipeRepository.NewRecipeRepository(dbConn)
	recipeUsecase := recipeUsecase.NewRecipeUsecase(recipeRepository)
	recipeHandler.NewRecipeHandler(g, authMiddleware, recipeUsecase)

	// set gin router with defined application port
	server := &http.Server{
		Addr:    ":" + application.Port,
		Handler: g,
	}
	// listen and serve
	go func() {
		log.Info("[Bootstrap] starting endeus")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[Bootstrap] error starting endeus: %v", err)
		}
	}()
	// prepare graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	// shutdown application
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("[Bootstrap] shutting down application: %v", err)
	}
	// close database connection
	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Fatalf("[Bootstrap] error closing database connection: %s", err.Error())
		}
		log.Info("[Bootstrap] database connection gracefully closed...")
	}()
	return nil
}
