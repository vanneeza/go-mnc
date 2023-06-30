package server

import (
	"log"
	"os"

	"github.com/vanneeza/go-mnc/api"
	"github.com/vanneeza/go-mnc/config"
	"github.com/vanneeza/go-mnc/utils/pkg"
)

func Run() error {

	db, err := config.InitDB()
	if err != nil {
		return err
	}
	defer config.CloseDB(db)
	jwtKey := os.Getenv("JWT_KEY")
	router := api.Run(db, jwtKey)
	serverAddress := pkg.GetEnv("SERVER_ADDRESS")
	log.Printf("Server is running on address %s\n", serverAddress)
	if err := router.Run(serverAddress); err != nil {
		return err
	}

	return nil
}
