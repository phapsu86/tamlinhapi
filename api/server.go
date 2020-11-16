package api

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/phapsu86/tamlinhapi/api/controllers"
	//"github.com/phapsu86/tamlinhapi/api/seed"
)

var server = controllers.Server{}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("sad .env file found")
	}
}

func Run() {

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	server.InitializeStorage(os.Getenv("STOREGRADE_ENDPOINT"),os.Getenv("STOREGRADE_ACCESSKEY_ID"),os.Getenv("STOREGRADE_SECRET_ACCESS_KEY"))
	//seed.Load(server.DB)

	server.Run(":8080")

}
