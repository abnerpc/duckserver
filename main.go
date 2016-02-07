package main

import (
	"log"
	"net/http"
)

func init() {
	// load config
	if err := LoadConfiguration(); err != nil {
		log.Fatal("Error loading the configuration file")
		return
	}
}

func main() {

	SetupAPIHandlers()
	SetupUploadHandler()

	static := http.FileServer(http.Dir("static"))
	http.Handle("/", UserSecureMiddleware(static))

	log.Println("Listening on port 8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
