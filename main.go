package main

import (
	"log"
	"net/http"
	"os"
)

func getHostPort() string {
	host := ""
	port := "7333"

	if envHost, exists := os.LookupEnv("HOST"); exists {
		host = envHost
	}

	if envPort, exists := os.LookupEnv("PORT"); exists {
		port = envPort
	}

	return host + ":" + port
}

func main() {
	addr := getHostPort()
	http.HandleFunc("/", YoutubeHandler)
	log.Printf("Listening on %s...\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
