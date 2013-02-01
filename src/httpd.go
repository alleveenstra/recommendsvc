package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"recommendsvc"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage %s [input-file.json] [listen-address]", os.Args[0])
		os.Exit(1)
	}
	filebytes, err := ioutil.ReadFile(os.Args[1])
	recommendsvc.Fatal_error(err)
	log.Println("Loading places...")
	var places []recommendsvc.Place
	json_err := json.Unmarshal(filebytes, &places)
	recommendsvc.Fatal_error(json_err)
	log.Println("Starting webserver...")
	http.HandleFunc("/geo", recommendsvc.Build_geo_handler(places))
	http.HandleFunc("/locality", recommendsvc.Build_locality_handler(places))
	http_err := http.ListenAndServe(os.Args[2], nil)
	recommendsvc.Fatal_error(http_err)
}
