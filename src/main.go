package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if candidate, found := os.LookupEnv("ERROR_RATE"); found {
		if rate, err := strconv.ParseFloat(candidate, 64); err == nil && rate >= rand.Float64() {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(http.StatusText(http.StatusForbidden)))
			return
		}
	}
	if revision, found := os.LookupEnv("K_REVISION"); found {
		fmt.Fprintf(w, "Revision: %s\n", revision)
		return
	}
	fmt.Fprintf(w, "Hello, world!\n")
}

func main() {
	http.HandleFunc("/", handler)

	port := "8080"
	if candidate, found := os.LookupEnv("PORT"); found {
		port = candidate
	}
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
