package main

import (
	"log"
	"net/http"
	"os"

	"webhook/router"

	"github.com/joho/godotenv"
)

func main() {
	initEnv()

	log.Println("Listening on localhost:8000...")
	http.ListenAndServe(getPort(), router.MakeHTTPHandler())

}
func getPort() string {
	p := os.Getenv("PORT")
	if p != "" {
		return ":" + p
	}
	return ":8000"
}
func initEnv() {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
